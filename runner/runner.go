package runner

import (
	"fmt"
	"strconv"
	"strings"

	d "github.com/fastbear1/quack/drivers"
	utils "github.com/fastbear1/quack/internal"
)

func Run(conf *utils.ConfigYaml) {
	// step 1: check connection to database
	var dbTablesMeta []d.TableMeta

	drv, err := d.GetDriver(conf.Database.Type)
	utils.CheckErrLite(err)

	// step 1.1: get avaiable tables from database
	dbTables, err := drv.GetTablesList(conf)
	fmt.Println(dbTables)
	utils.CheckErrLite(err)
	for _, tableName := range dbTables {
		dbTablesMeta = append(dbTablesMeta,
			d.TableMeta{
				Name: tableName,
			},
		)
	}
	// step 1.2: get table columns meta data
	for i := 0; i < len(dbTablesMeta); i++ {
		res, err := drv.GetTableColumnsMeta(conf, dbTablesMeta[i].Name)
		utils.CheckErrLite(err)
		dbTablesMeta[i].Columns = res
		// step 1.3: get table indices information
		idx, err := drv.GetTableIndices(conf, dbTablesMeta[i].Name)
		utils.CheckErrLite(err)
		dbTablesMeta[i].Indeces = idx
		// step 1.4: get table references
		ref, err := drv.GetTableReferences(conf, dbTablesMeta[i].Name)
		utils.CheckErrLite(err)
		dbTablesMeta[i].References = ref
	}

	// step 2: Scan models directory for gorm struct definitions
	var gormStructMeta []d.TableMeta

	StructRaw, err := Scan(conf)
	utils.CheckErrLite(err)

	for i := 0; i < len(StructRaw); i++ {
		// step 2.1: get gorm tables metadata
		gsmeta := parseModelStruct(StructRaw[i], drv)
		// step 2.2: parse gorm index information
		for _, f := range StructRaw[i].Fields {
			if index, ok := parseIndicesTag(gsmeta.Name, drv.TransformName(f.FieldName), f.FieldTag); ok {
				gsmeta.Indeces = append(gsmeta.Indeces, index)
			}
		}
		// step 2.3: parse FK and embed strcuture
		for _, f := range StructRaw[i].ReferenceFields {
			reference := parseReferenceEmbedStructs(drv, StructRaw[i].Name, f.FieldType, f.FieldTag)
			gsmeta.References = append(gsmeta.References, reference)
		}
		gormStructMeta = append(gormStructMeta, gsmeta)
	}

	// step3: Compare current state of metadata for database tables and gorm structures
	funcList, err := compareMetaState(dbTablesMeta, gormStructMeta)
	utils.CheckErrLite(err)
	var sqlUp, sqlDown []string
	for _, f := range funcList {
		up, down := f(drv)
		sqlUp = append(sqlUp, up)
		sqlDown = append(sqlDown, down)
	}

	// step 4: Write sql- Up and Down commands to file
	if len(sqlUp) != 0 || len(sqlDown) != 0 {
		writeToFile(conf, sqlUp, sqlDown)
	} else {
		fmt.Println("Gorm struct and DB tables already synchronized")
	}
}

func parseModelStruct(data ModelStruct, drv d.DbHandler) d.TableMeta {
	model := d.TableMeta{
		Name:       drv.TransformName(data.Name),
		Columns:    make([]d.Column, 0),
		Indeces:    make([]d.IndexMeta, 0),
		References: make([]d.ReferenceMeta, 0),
	}

	for _, f := range data.Fields {
		column := d.Column{
			TableName:         model.Name,
			ColumnName:        drv.TransformName(f.FieldName),
			DataType:          drv.TransformType(f.FieldType),
			IsNullable:        false,
			ColumnDefault:     "",
			IsPrimary:         false,
			PrimaryConstraint: "",
		}
		// check that we got non empty and correct tag value
		if f.FieldTag != `` {
			if strings.HasPrefix(f.FieldTag, "gorm:") {
				prefix := "gorm:"
				tag := strings.TrimPrefix(f.FieldTag, prefix)
				tag = tag[1 : len(tag)-1]
				parseTag(&column, tag)
			} else {
				fmt.Printf("Not found any tag for field %s in struct %s\n", f.FieldName, data.Name)
			}
		}
		if column.ColumnDefault != "" {
			column.ColumnDefault = drv.TransformDefault(column.DataType, column.ColumnDefault)
		}
		model.Columns = append(model.Columns, column)
	}
	return model
}

func parseTag(col *d.Column, tag string) {
	for _, t := range strings.Split(tag, ";") {
		if strings.Contains(t, ":") {
			splitval := strings.Split(t, ":")
			key := splitval[0]
			val := splitval[1]
			switch key {
			case "type":
				col.DataType = val
			case "default":
				col.ColumnDefault = val
			}
		} else {
			switch t {
			case "primary_key":
				col.IsPrimary = true
			case "not null":
				col.IsNullable = false
			case "null":
				col.IsNullable = true
			}
		}
	}
}

func parseIndicesTag(table string, column string, tag string) (d.IndexMeta, bool) {
	var (
		idxmeta  d.IndexMeta = d.IndexMeta{}
		idxfound bool        = false
	)
	tag = strings.TrimPrefix(tag, "gorm:")
	if tag == "" {
		return idxmeta, idxfound
	}
	tag = tag[1 : len(tag)-1]

	for _, value := range strings.Split(tag, ";") {
		if value != "" {
			if strings.Contains(value, "primary_key") {
				// Skip primary keys indices
				continue
			}
			v := strings.Split(value, ":")
			k := strings.TrimSpace(v[0])
			if k == "index" || k == "uniqueindex" {
				var (
					name       string
					tag        = strings.Join(v[1:], ":")
					idx        = strings.IndexByte(tag, ',')
					tagSetting = strings.Join(strings.Split(tag, ",")[1:], ",")
					settings   = ParseTagSetting(tagSetting, ",")
				)

				if idx == -1 {
					idx = len(tag)
				}

				name = tag[0:idx]
				if name == "" {
					name = createIndexName(table, []string{column}, settings["EXPRESSION"])
				}
				uniqidx := false
				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					uniqidx = true
				}
				priority, err := strconv.Atoi(settings["PRIORITY"])
				if err != nil {
					priority = 10
				}
				if settings["TYPE"] == "" {
					//TODO: set drier default index type
					settings["TYPE"] = "btree"
				}
				idxmeta = d.IndexMeta{
					TableName: table,
					Name:      name,
					Unique:    uniqidx,
					Type:      settings["TYPE"],
					Where:     settings["WHERE"],
					Option:    settings["OPTION"],
					Parsed:    true,
					Columns: []d.IndexOption{{
						Field:      column,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Priority:   priority,
					}},
				}
				idxfound = true
			}
		}
	}
	return idxmeta, idxfound
}

func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	var parsedNames []string
	for i := 0; i < len(names); i++ {
		s := names[i]
		for strings.HasSuffix(s, "\\") && i+1 < len(names) {
			i++
			s = s[:len(s)-1] + sep + names[i]
		}
		parsedNames = append(parsedNames, s)
	}

	for _, tag := range parsedNames {
		values := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))
		if len(values) >= 2 {
			val := strings.Join(values[1:], ":")
			val = strings.ReplaceAll(val, `\"`, `"`)
			settings[k] = val
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}

func createIndexName(table string, columns []string, exp string) string {
	indexName := "idx_"
	columnsSuffix := ""
	expSuffix := ""
	if exp != "" {
		expParts := strings.Split(exp, "(")
		expSuffix += fmt.Sprintf("_%s", expParts[0])
	}
	for _, c := range columns {
		columnsSuffix += fmt.Sprintf("_%s", c)
	}
	indexName += table + columnsSuffix + expSuffix
	return indexName
}

func parseReferenceEmbedStructs(drv d.DbHandler, table string, reftable string, tag string) d.ReferenceMeta {
	// Example: gorm:"foreignKey:UserName;references:Name;referenceName:fk_auth_users_users;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;
	var ref = d.ReferenceMeta{
		TableName: table,
	}
	tag = strings.TrimPrefix(tag, "gorm:")
	tag = tag[1 : len(tag)-1]
	for _, value := range strings.Split(tag, ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(v[0])
			switch k {
			case "referenceName":
				ref.Name = v[1]
			case "foreignKey":
				ref.Column = drv.TransformName(v[1])
			case "references":
				ref.RefColumn = drv.TransformName(v[1])
			case "constraint":
				ref.RefOptions = ""
				constr := strings.Join(v[1:], ":")
				constParts := strings.Split(constr, ",")
				for i, c := range constParts {
					if i > 0 {
						ref.RefOptions += " "
					}
					action := strings.Split(c, ":")
					constrAction := transformAction(string(action[0]))
					ref.RefOptions += strings.ToUpper(constrAction) + " " + strings.ToUpper(string(action[1]))
				}
			}
		}
	}
	ref.RefTable = drv.TransformName(reftable)
	if ref.RefColumn == "" {
		ref.RefColumn = "id"
	}
	if ref.Name == "" {
		//fk_TableName_ReferenceTableName_ForeignKeyColumn_ReferenceColumn
		ref.Name = "fk_" + drv.TransformName(table) + "_" +
			drv.TransformName(reftable) + "_" +
			ref.Column + "_" +
			drv.TransformName(ref.RefColumn)
	}
	return ref
}

func transformAction(action string) string {
	// TODO: driver depend function, delegate transformation to driver
	defaction := "ON DELETE"
	if action == "OnUpdate" {
		defaction = "ON UPDATE"
	}
	return defaction
}

func compareMetaState(dbmeta []d.TableMeta, gmeta []d.TableMeta) ([]func(drv d.DbHandler) (string, string), error) {
	var (
		funcList []func(drv d.DbHandler) (string, string)
		dbmap    map[string]d.TableMeta
		gmap     map[string]d.TableMeta
	)

	if len(dbmeta) == 0 {
		// return create table for all objects in gmeta
		for _, str := range gmeta {
			funcList = append(funcList, str.CreateTable)
			for _, i := range str.Indeces {
				funcList = append(funcList, i.CreateIndex)
			}
		}
		// job done, not enough data for comparing
		return funcList, nil
	}
	var (
		left  = make([]string, 0)
		right = make([]string, 0)
	)
	for _, l := range gmeta {
		left = append(left, l.Name)
		gmap[l.Name] = l
	}
	for _, r := range dbmeta {
		right = append(right, r.Name)
		dbmap[r.Name] = r
	}

	toDelete, toCreate := getCatalogData(left, right)

	if len(toDelete) > 0 || len(toCreate) > 0 {
		//TODO: unefficient
		for _, cr := range toCreate {
			for _, l := range gmeta {
				if cr == l.Name {
					funcList = append(funcList, l.CreateTable)
					for _, idx := range l.Indeces {
						funcList = append(funcList, idx.CreateIndex)
					}
				}
			}
		}

		for _, dt := range toDelete {
			for _, r := range dbmeta {
				if dt == r.Name {
					funcList = append(funcList, r.DeleteTable)
				}
			}
		}
	}

	// Check columns
	// check all column are exists
	for name, gtable := range gmap {
		if dbtable, ok := dbmap[name]; !ok {
			// Skipping tables that are not exists for now
			continue
		} else {
			toCreateCol, toDeleteCol := StateDifference(gtable.Columns, dbtable.Columns)
			for _, c := range toCreateCol {
				//col := c.(d.Column)
				funcList = append(funcList, c.CreateColumn)
			}
			for _, c := range toDeleteCol {
				//col := c.(d.Column)
				funcList = append(funcList, c.DeleteColumn)
			}
		}
	}

	// Not implemented
	return funcList, nil
}

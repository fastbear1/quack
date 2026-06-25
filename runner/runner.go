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
			index := parseIndicesTag(gsmeta.Name, f.FieldName, f.FieldTag)
			gsmeta.Indeces = append(gsmeta.Indeces, index)
		}
		// step 2.3: parse FK and embed strcuture
		for _, f := range StructRaw[i].Fields {
			reference := parseReferenceEmbedStructs(f.FieldName)
			gsmeta.References = append(gsmeta.References, reference)
		}
		gormStructMeta = append(gormStructMeta, gsmeta)
	}
	//fmt.Printf("%+v\n", gormStructMeta)

	// step3: Compare current state of metadata for database tables and gorm structures
	funcList, err := compareMetaState(dbTablesMeta, gormStructMeta)
	utils.CheckErrLite(err)
	var sqlUp, sqlDown []string
	for _, f := range funcList {
		up, down := f(conf, drv)
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

func parseIndicesTag(table string, column string, tag string) d.IndexMeta {
	var idxmeta d.IndexMeta = d.IndexMeta{}
	for _, value := range strings.Split(tag, ";") {
		if value != "" {
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
				fmt.Println(tag)
				fmt.Println(idx)
				fmt.Println(tagSetting)

				if idx == -1 {
					idx = len(tag)
				}

				name = tag[0:idx]
				if name == "" {
					name = fmt.Sprintf("idx_%s_%s", table, column)
				}

				uniqidx := false
				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					uniqidx = true
				}

				priority, err := strconv.Atoi(settings["PRIORITY"])
				if err != nil {
					priority = 10
				}

				idxmeta = d.IndexMeta{
					Name:   name,
					Unique: uniqidx,
					Type:   settings["TYPE"],
					Where:  settings["WHERE"],
					Option: settings["OPTION"],
					Parsed: true,
					Columns: []d.IndexOption{{
						Field:      column,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Priority:   priority,
					}},
				}
			}
		}
	}
	return idxmeta
}

func parseReferenceEmbedStructs(name string) d.ReferenceMeta {
	return d.ReferenceMeta{}
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

func compareMetaState(dbmeta []d.TableMeta, gmeta []d.TableMeta) ([]func(conf *utils.ConfigYaml, drv d.DbHandler) (string, string), error) {
	var funcList []func(conf *utils.ConfigYaml, drv d.DbHandler) (string, string)
	if len(dbmeta) == 0 {
		// return create table for all objects in gmeta
		for _, str := range gmeta {
			funcList = append(funcList, str.CreateTable)
		}
	}
	// Not implemented
	//var metamap map[string]TableMeta
	//for i := 0; i < len(gmeta); i++ {
	//	metamap[(gmeta)[i].Name] = &gmeta[i]
	//}
	return funcList, nil
}

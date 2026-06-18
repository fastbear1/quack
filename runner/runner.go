package runner

import (
	"fmt"
	"strings"

	d "github.com/fastbear1/quack/drivers"
	utils "github.com/fastbear1/quack/internal"
)

func Run(conf *utils.ConfigYaml) {
	// step 1: check connection to database
	var dbTablesMeta []d.TableMeta

	drv, err := d.GetDriver(conf.Database.Type)
	utils.CheckErrLite(err)

	dbTables, err := drv.GetTablesList(conf)
	utils.CheckErrLite(err)
	for _, tableName := range dbTables {
		dbTablesMeta = append(dbTablesMeta,
			d.TableMeta{
				Name: tableName,
			},
		)
	}

	for i := 0; i < len(dbTablesMeta); i++ {
		res, err := drv.GetTableColumnsMeta(conf, dbTablesMeta[i].Name)
		utils.CheckErrLite(err)
		dbTablesMeta[i].Columns = res
	}

	// step 2: Scan models directory for gorm struct definitions
	var gormStructMeta []d.TableMeta

	StructRaw, err := Scan(conf)
	utils.CheckErrLite(err)
	gormStructMeta = parseModelStruct(StructRaw, drv)

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
	if len(sqlUp) != 0 && len(sqlDown) != 0 {
		writeToFile(conf, sqlUp, sqlDown)
	} else {
		fmt.Println("Gorm struct and DB tables already synchronized")
	}
}

func parseModelStruct(data []ModelStruct, drv d.DbHandler) []d.TableMeta {
	var allModels []d.TableMeta
	for _, m := range data {
		model := d.TableMeta{
			Name:       drv.TransformName(m.Name),
			Columns:    make([]d.Column, 0),
			Indeces:    make([]d.IndexMeta, 0),
			References: make([]d.ReferenceMeta, 0),
		}

		for _, f := range m.Fields {
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
					fmt.Printf("Not found any tag for field %s in struct %s\n", f.FieldName, m.Name)
				}
			}
			model.Columns = append(model.Columns, column)
		}
		allModels = append(allModels, model)
	}
	return allModels
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

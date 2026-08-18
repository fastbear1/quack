package runner

import (
	"strconv"
	"strings"

	. "github.com/fastbear1/quack/schema"
)

func parseModelStruct(drv DbInterface, data ModelStruct) TableMeta {
	model := TableMeta{
		Name:       drv.TransformName(data.Name),
		Columns:    make([]Column, 0),
		Indeces:    make([]IndexMeta, 0),
		References: make([]ReferenceMeta, 0),
	}

	for _, field := range data.Fields {
		parseModelTags(drv, &model, field)
	}
	return model
}

func parseModelTags(drv DbInterface, model *TableMeta, field FieldStruct) {
	var idx IndexMeta
	// Default column declaration
	column := Column{
		TableName:  model.Name,
		ColumnName: drv.TransformName(field.FieldName),
		DataType:   drv.TransformType(field.FieldType),
		IsNullable: func(force bool) bool {
			if force {
				return true
			}
			return false
		}(field.ForceNull),
		ColumnDefault:  "",
		IsPrimary:      false,
		PrimaryOptions: PrimaryOptions{},
	}

	if field.FieldTag != `` && strings.HasPrefix(field.FieldTag, "gorm:") {
		tag := strings.TrimPrefix(field.FieldTag, "gorm:")
		tag = tag[1 : len(tag)-1] // remove parent double quotes
		for _, value := range strings.Split(tag, ";") {
			kit := strings.Split(value, ":")
			kitLen := len(kit)
			k := strings.TrimSpace(kit[0])
			switch k {
			case "type":
				if val := getValue(kit, kitLen, 1); val != "" {
					column.DataType = val
				}
			case "default":
				if val := getValue(kit, kitLen, 1); val != "" {
					column.ColumnDefault = drv.TransformDefault(column.DataType, val)
				}
			case "primaryKey":
				column.IsPrimary = true
			case "not null":
				// already initialized with false
				if !field.ForceNull {
					column.IsNullable = false
				}
			case "null":
				column.IsNullable = true
			case "index", "uniqueindex":
				idx.TableName = model.Name
				var (
					name       string
					tag        = strings.Join(kit[1:], ":")
					nameId     = strings.IndexByte(tag, ',')
					tagSetting = strings.Join(strings.Split(tag, ",")[1:], ",")
					settings   = ParseTagSetting(tagSetting, ",")
				)

				if nameId == -1 {
					nameId = len(tag)
				}
				name = tag[0:nameId]
				if name == "" {
					name = drv.CreateIndexName(model.Name, []string{column.ColumnName}, settings["expression"])
				}
				uniqidx := false
				if (k == "uniqueindex") || settings["unique"] != "" {
					uniqidx = true
				}
				priority, err := strconv.Atoi(settings["priority"])
				if err != nil {
					priority = 10
				}
				if settings["type"] == "" {
					// set driver default index type
					settings["type"] = "btree"
				}

				idxFound := false
				// check if index already exists (composite index)
				for n, i := range model.Indeces {
					if i.Name == name {
						idxFound = true
						model.Indeces[n].Columns = append(model.Indeces[n].Columns,
							IndexOption{
								Field:      column.ColumnName,
								Expression: settings["expression"],
								Sort:       settings["sort"],
								Collate:    settings["collate"],
								Priority:   priority,
							})

					}
				}

				if !idxFound {
					idx.Name = name
					idx.Unique = uniqidx
					idx.Type = settings["type"]
					idx.Where = settings["where"]
					idx.Option = settings["option"]
					idx.Parsed = true
					idx.Columns = []IndexOption{{
						Field:      column.ColumnName,
						Expression: settings["expression"],
						Sort:       settings["sort"],
						Collate:    settings["collate"],
						Priority:   priority,
					}}
					model.Indeces = append(model.Indeces, idx)
				}
			}
		}
	}
	// set primary key status
	// postgres specification only
	if column.IsPrimary {
		if column.DataType == "serial" {
			column.DataType = drv.TransformType(column.DataType)
			column.PrimaryOptions.IsSerial = true
		}
		if strings.HasPrefix(column.ColumnDefault, "generated") {
			column.PrimaryOptions.IsIdentity = true
			column.ColumnDefault = strings.ToUpper(column.ColumnDefault)
		}
	}
	model.Columns = append(model.Columns, column)
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
		k := strings.TrimSpace(values[0])
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

func getValue(data []string, dl int, n int) string {
	if n < dl {
		return data[n]
	}
	return ""
}

func parseReferenceEmbedStructs(drv DbInterface, table string, reftable string, tag string) ReferenceMeta {
	var ref = ReferenceMeta{
		TableName: drv.TransformName(table),
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
					constrAction := drv.TransformConstraintAction(string(action[0]))
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
		ref.Name = drv.CreateConstraintName(table, ref.Column, ref.RefTable, ref.RefColumn)
	}
	return ref
}

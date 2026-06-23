package runner

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	utils "github.com/fastbear1/quack/internal"
)

func formatNode(fset *token.FileSet, node ast.Node) (string, error) {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TODO: return error not panic
func mustFormatNode(fset *token.FileSet, node ast.Node) string {
	str, err := formatNode(fset, node)
	utils.CheckErrLite(err)
	return str
}

type FieldStruct struct {
	FieldName string
	FieldType string
	FieldTag  string
}

type ModelStruct struct {
	Name   string
	Fields []FieldStruct
}

func getStructs(conf *utils.ConfigYaml, fset *token.FileSet, file *ast.File) []ModelStruct {
	var structdef []ModelStruct
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				tspec := spec.(*ast.TypeSpec)
				if structType, ok := tspec.Type.(*ast.StructType); ok {
					structName := tspec.Name.Name
					if utils.InArray(conf.Models.Exclude, structName) {
						continue
					}
					modelData := ModelStruct{
						Name: structName,
					}
					var fieldData FieldStruct

					for _, field := range structType.Fields.List {
						if len(field.Names) > 0 {
							fieldData.FieldName = field.Names[0].String()
						} else {
							fieldData.FieldName = "Embed field"
						}
						fieldData.FieldType = mustFormatNode(fset, field.Type)
						if field.Tag != nil {
							fieldData.FieldTag = field.Tag.Value[1 : len(field.Tag.Value)-1]
						}
						modelData.Fields = append(modelData.Fields, fieldData)
					}
					structdef = append(structdef, modelData)
				}
			}
		}
	}
	return structdef
}

func visitFiles(conf *utils.ConfigYaml, fset *token.FileSet, files map[string]*ast.File) ([]ModelStruct, error) {
	var models []ModelStruct
	for _, file := range files {
		//fmt.Printf("/* file: %s */\n", filename)
		var structList []ModelStruct
		structList = getStructs(conf, fset, file)
		models = append(models, structList...)
	}
	return models, nil
}

func Scan(conf *utils.ConfigYaml) ([]ModelStruct, error) {
	var path string = fmt.Sprintf("./%s", conf.Models.Path)

	// Create the AST file set.
	fset := token.NewFileSet()

	// Parse all files in the directory.
	pkgs, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) == ".go"
	}, parser.ParseComments)

	//pkgs, err := parser.ParseFile(fset, "", path, parser.Mode())

	if err != nil {
		fmt.Println("Error parsing directory:", err)
		return []ModelStruct{}, nil
	}

	// Map to hold package files.
	files := make(map[string]*ast.File)
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			files[filepath.Join(path, name)] = file
		}
	}
	// Visit all files and collect struct definitions.
	return visitFiles(conf, fset, files)
}

/*
func parseFieldIndexes(field *Field) (indexes []Index, err error) {
	for _, value := range strings.Split(field.Tag.Get("gorm"), ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "INDEX" || k == "UNIQUEINDEX" {
				var (
					name       string
					tag        = strings.Join(v[1:], ":")
					idx        = strings.IndexByte(tag, ',')
					tagSetting = strings.Join(strings.Split(tag, ",")[1:], ",")
					settings   = ParseTagSetting(tagSetting, ",")
					length, _  = strconv.Atoi(settings["LENGTH"])
				)

				if idx == -1 {
					idx = len(tag)
				}

				name = tag[0:idx]
				if name == "" {
					subName := field.Name
					const key = "COMPOSITE"
					if composite, found := settings[key]; found {
						if len(composite) == 0 || composite == key {
							err = fmt.Errorf(
								"the composite tag of %s.%s cannot be empty",
								field.Schema.Name,
								field.Name)
							return
						}
						subName = composite
					}
					name = field.Schema.namer.IndexName(
						field.Schema.Table, subName)
				}

				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					settings["CLASS"] = "UNIQUE"
				}

				priority, err := strconv.Atoi(settings["PRIORITY"])
				if err != nil {
					priority = 10
				}

				indexes = append(indexes, Index{
					Name:    name,
					Class:   settings["CLASS"],
					Type:    settings["TYPE"],
					Where:   settings["WHERE"],
					Comment: settings["COMMENT"],
					Option:  settings["OPTION"],
					Fields: []IndexOption{{
						Field:      field,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Length:     length,
						Priority:   priority,
					}},
				})
			}
		}
	}

	err = nil
	return
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
*/

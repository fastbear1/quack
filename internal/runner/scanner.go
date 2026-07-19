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
	Name            string
	Fields          []FieldStruct
	EmbedFields     []string
	ReferenceFields []FieldStruct
}

func collectEmbedFileds(data []ModelStruct, cache map[string]ModelStruct) []ModelStruct {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i].EmbedFields); j++ {
			if model, exist := cache[data[i].EmbedFields[j]]; exist {
				data[i].Fields = append(model.Fields, data[i].Fields...)
				//data[i].EmbedFields = []string{}
			}
		}
	}
	return data
}

func getStructs(conf *utils.ConfigYaml, fset *token.FileSet, file *ast.File) ([]ModelStruct, map[string]ModelStruct) {
	// TODO: refactor this function
	// TODO: cached structs should be a map
	var (
		structdef    []ModelStruct
		cachedStruct = make(map[string]ModelStruct)
	)
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				tspec := spec.(*ast.TypeSpec)
				if structType, ok := tspec.Type.(*ast.StructType); ok {
					structName := tspec.Name.Name
					modelData := ModelStruct{
						Name:            structName,
						Fields:          make([]FieldStruct, 0),
						EmbedFields:     make([]string, 0),
						ReferenceFields: make([]FieldStruct, 0),
					}
					var fieldData FieldStruct

					for _, field := range structType.Fields.List {
						if len(field.Names) > 0 {
							fieldData.FieldName = field.Names[0].String()

							if r, ok := field.Type.(*ast.Ident); ok && r.Obj != nil {
								// FK column declaration
								var refField = FieldStruct{}
								refField.FieldName = field.Names[0].String()
								refField.FieldType = r.Name
								if field.Tag != nil {
									refField.FieldTag = field.Tag.Value[1 : len(field.Tag.Value)-1]
								}
								modelData.ReferenceFields = append(modelData.ReferenceFields, refField)
							} else {
								fieldData.FieldType = mustFormatNode(fset, field.Type)
								if field.Tag != nil {
									fieldData.FieldTag = field.Tag.Value[1 : len(field.Tag.Value)-1]
								}
								modelData.Fields = append(modelData.Fields, fieldData)
							}
						} else {
							// Anonymous embed field
							if r, ok := field.Type.(*ast.Ident); ok {
								modelData.EmbedFields = append(modelData.EmbedFields, r.Name)
							}
							continue
						}
					}

					if utils.InArray(conf.Models.Exclude, structName) {
						cachedStruct[structName] = modelData
					} else {
						structdef = append(structdef, modelData)
					}
				}
			}
		}
	}
	return structdef, cachedStruct
}

func visitFiles(conf *utils.ConfigYaml, fset *token.FileSet, files map[string]*ast.File) ([]ModelStruct, error) {
	var models []ModelStruct
	var cached = make(map[string]ModelStruct)
	for _, file := range files {
		var structList []ModelStruct
		structList, cacheList := getStructs(conf, fset, file)
		models = append(models, structList...)
		for k, v := range cacheList {
			cached[k] = v
		}
	}
	models = collectEmbedFileds(models, cached)
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
		return []ModelStruct{}, err
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

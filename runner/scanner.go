package runner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"

	utils "github.com/fastbear1/quack/internal"
)

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

func FormatNode(fset *token.FileSet, node ast.Node) string {
	var fType string
	switch t := node.(type) {
	case *ast.Ident:
		fType = t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			fType = x.String() + "." + t.Sel.String()
		}
	}
	return fType
}

func collectEmbedFileds(data []ModelStruct, cache map[string]ModelStruct) []ModelStruct {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i].EmbedFields); j++ {
			if model, exist := cache[data[i].EmbedFields[j]]; exist {
				data[i].Fields = append(model.Fields, data[i].Fields...)
			}
		}
	}
	return data
}

func getStructs(conf *utils.ConfigYaml, fset *token.FileSet, file *ast.File) ([]ModelStruct, []ModelStruct) {
	var (
		structdef    []ModelStruct
		cachedStruct []ModelStruct
	)
	ast.Inspect(file, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if e, ok := t.Type.(*ast.StructType); ok {
				structName := t.Name.Name
				modelData := ModelStruct{
					Name:            structName,
					Fields:          make([]FieldStruct, 0),
					EmbedFields:     make([]string, 0),
					ReferenceFields: make([]FieldStruct, 0),
				}
				var fieldData FieldStruct

				for _, field := range e.Fields.List {
					if len(field.Names) > 0 {
						fieldData.FieldName = field.Names[0].String()

						if r, ok := field.Type.(*ast.Ident); ok && r.Obj != nil {
							// FK column declaration
							var refField = FieldStruct{}
							refField.FieldName = field.Names[0].String()
							refField.FieldType = FormatNode(fset, field.Type)
							if field.Tag != nil {
								refField.FieldTag = field.Tag.Value[1 : len(field.Tag.Value)-1]
							}
							modelData.ReferenceFields = append(modelData.ReferenceFields, refField)
						} else {
							fieldData.FieldType = FormatNode(fset, field.Type)
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
					}
				}

				if utils.InArray(conf.Models.Exclude, structName) {
					cachedStruct = append(cachedStruct, modelData)
				} else {
					structdef = append(structdef, modelData)
				}

			}
		}
		return true
	})
	return structdef, cachedStruct
}

func visitFiles(conf *utils.ConfigYaml, fset *token.FileSet, files []*ast.File) ([]ModelStruct, error) {
	var models []ModelStruct
	var cached = make(map[string]ModelStruct)
	for _, file := range files {
		var structList []ModelStruct
		structList, cacheList := getStructs(conf, fset, file)
		models = append(models, structList...)
		for _, model := range cacheList {
			cached[model.Name] = model
		}
	}
	models = collectEmbedFileds(models, cached)
	return models, nil
}

func Scan(conf *utils.ConfigYaml) ([]ModelStruct, error) {
	pFiles := make([]*ast.File, 0)
	var pDir string = fmt.Sprintf("./%s", conf.Models.Path)
	// Create the AST file set.
	fset := token.NewFileSet()

	// read file names in directory
	files, err := os.ReadDir(pDir)
	if err != nil {
		fmt.Println("Error parsing directory:", err)
		return nil, err
	}

	// Parse all files in direcotry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fn := path.Join(pDir, file.Name())
		// parse file and generate ast tree
		ftree, err := parser.ParseFile(fset, fn, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Error parsing file: ", file.Name())
		}
		pFiles = append(pFiles, ftree)
	}

	// Visit all files and collect struct definitions.
	return visitFiles(conf, fset, pFiles)
}

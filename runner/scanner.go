package runner

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fastbear1/quack/drivers"
	utils "github.com/fastbear1/quack/internal"
)

func getDbMeta(conf *utils.ConfigYaml) ([]string, error) {
	drv, err := drivers.GetDriver(conf.Database.Type)
	fmt.Println(drv)
	utils.CheckErrLite(err)
	data, err := drv.GetData(conf)
	fmt.Println(data)
	return data, err
}

func formatNode(fset *token.FileSet, node ast.Node) (string, error) {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func mustFormatNode(fset *token.FileSet, node ast.Node) string {
	str, err := formatNode(fset, node)
	if err != nil {
		panic(err)
	}
	return str
}

func printComments(fset *token.FileSet, file *ast.File, pos token.Pos, indent int) {
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if fset.Position(comment.Pos()).Line == fset.Position(pos).Line-1 {
				fmt.Printf("%s%s\n", strings.Repeat(" ", indent), comment.Text)
			}
		}
	}
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

func printStructs(fset *token.FileSet, file *ast.File) []ModelStruct {
	var structdef []ModelStruct
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				tspec := spec.(*ast.TypeSpec)
				if structType, ok := tspec.Type.(*ast.StructType); ok {
					// fmt.Println(structType)
					// fmt.Println(structType.Fields.NumFields())
					structName := tspec.Name.Name
					modelData := ModelStruct{
						Name: structName,
					}
					var fieldData FieldStruct

					// Print comments associated with the struct
					// printComments(fset, file, tspec.Pos(), 0)
					fmt.Printf("type %s struct {\n", structName)
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

						if field.Names != nil {
							for _, name := range field.Names {
								// Print comments associated with the field
								printComments(fset, file, name.Pos(), 4)
								fmt.Printf("    %s %s", name.Name, mustFormatNode(fset, field.Type))
								if field.Tag != nil {
									fmt.Printf(" %s", field.Tag.Value)
								}
								fmt.Println("")
							}
						} else { // anonymous field
							// Print comments associated with the field
							printComments(fset, file, field.Type.Pos(), 4)
							fmt.Printf("    %s", mustFormatNode(fset, field.Type))
							if field.Tag != nil {
								fmt.Printf(" %s", field.Tag.Value)
							}
							fmt.Println("")
						}
						modelData.Fields = append(modelData.Fields, fieldData)
					}
					structdef = append(structdef, modelData)
					fmt.Println("}")
				}
			}
		}
	}
	return structdef
}

func visitFiles(fset *token.FileSet, files map[string]*ast.File) []ModelStruct {
	var models []ModelStruct
	for filename, file := range files {
		fmt.Printf("/* file: %s */\n", filename)
		var structList []ModelStruct
		structList = printStructs(fset, file)
		models = append(models, structList...)
	}
	return models
}

func Scan() []ModelStruct {
	var path string
	flag.StringVar(&path, "path", "./models/", "Directory to parse")
	flag.Parse()

	// Create the AST file set.
	fset := token.NewFileSet()

	// Parse all files in the directory.
	pkgs, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
		//log.Infof("# file: %s", info.Name())
		return !info.IsDir() && filepath.Ext(info.Name()) == ".go"
	}, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing directory:", err)
		os.Exit(1)
	}

	// Map to hold package files.
	files := make(map[string]*ast.File)
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			files[filepath.Join(path, name)] = file
		}
	}
	fmt.Println(files)
	// Visit all files and collect struct definitions.
	return visitFiles(fset, files)
}

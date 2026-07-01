package drivers

import (
	"errors"

	utils "github.com/fastbear1/quack/internal"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type DbHandler interface {
	GetTablesList(conf *utils.ConfigYaml) ([]string, error)
	GetTableColumnsMeta(conf *utils.ConfigYaml, name string) ([]Column, error)
	GetTableIndices(conf *utils.ConfigYaml, name string) ([]IndexMeta, error)
	GetTableReferences(conf *utils.ConfigYaml, name string) ([]ReferenceMeta, error)
	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(col_type string, val string) string
	CreateTableStatement(table *TableMeta) (string, string)
	DropTableStatement(table *TableMeta) (string, string)
	CreateColumnStatement(col *Column) (string, string)
	DropColumnStatement(col *Column) (string, string)
	CreateIndexStatement(idx *IndexMeta) (string, string)
	DropIndexStatement(idx *IndexMeta) (string, string)
	CreateConstraintStatement(ref *ReferenceMeta) (string, string)
	DropConstraintStatement(ref *ReferenceMeta) (string, string)
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
	case "postgres":
		return &PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

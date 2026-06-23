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
	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(val string) string
	CreateTableStatement(conf *utils.ConfigYaml, table *TableMeta) (string, string)
	CreateIndexStatement(conf *utils.ConfigYaml, idx *IndexMeta) (string, string)
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
	case "postgres":
		return &PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

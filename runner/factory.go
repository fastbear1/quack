package runner

import (
	"errors"

	pg "github.com/fastbear1/quack/drivers/postgres"
	utils "github.com/fastbear1/quack/internal"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type Column struct {
	ColumnName        string
	DataType          string
	IsNullable        bool
	ColumnDefault     string
	IsPrimary         bool
	PrimaryConstraint string
}

type ReferenceMeta struct {
	Name           string
	RefColumn      string
	RefTable       string
	RefTableColumn string
	RefConstraint  string
}

type IndexMeta struct {
	Name        string
	IndexColumn string
	IndexType   string
}

type TableMeta struct {
	Name       string
	Columns    []Column
	References []ReferenceMeta
	Indeces    []IndexMeta
}

type DbHandler interface {
	GetTablesList(conf *utils.ConfigYaml) ([]string, error)
	GetTableColumnsMeta(conf *utils.ConfigYaml, name string) ([]Column, error)
	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(val string) string
	CreateTableStatement(conf *utils.ConfigYaml, table *TableMeta) (string, string)
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
	case "postgres":
		return &pg.PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

func (table *TableMeta) CreateTable(conf *utils.ConfigYaml, drv DbHandler) (string, string) {
	return drv.CreateTableStatement(conf, table)
}

package drivers

import (
	"context"
	"errors"

	utils "github.com/fastbear1/quack/internal"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type DbHandler interface {
	GetTablesList(ctx context.Context, conf *utils.ConfigYaml) ([]string, error)
	GetTableColumnsMeta(ctx context.Context, conf *utils.ConfigYaml, name string) ([]Column, error)
	GetTableIndices(ctx context.Context, conf *utils.ConfigYaml, name string) ([]IndexMeta, error)
	GetTableReferences(ctx context.Context, conf *utils.ConfigYaml, name string) ([]ReferenceMeta, error)
	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(col_type string, val string) string
	CreateTableStatement(table *TableMeta) string
	DropTableStatement(table *TableMeta) string
	CreateColumnStatement(col *Column) string
	AlterColumnStatement(col *Column) string
	DropColumnStatement(col *Column) string
	CreateIndexStatement(idx *IndexMeta) string
	DropIndexStatement(idx *IndexMeta) string
	CreateConstraintStatement(ref *ReferenceMeta) string
	DropConstraintStatement(ref *ReferenceMeta) string
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
	case "postgres":
		return &PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

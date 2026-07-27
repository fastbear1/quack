package runner

import (
	"context"
	"errors"

	pg "github.com/fastbear1/quack/drivers/pg_driver"
	utils "github.com/fastbear1/quack/internal"
	sch "github.com/fastbear1/quack/schema"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type OffDbInterface interface {
	GetTablesList(ctx context.Context, conf *utils.ConfigYaml) ([]string, error)
	GetTableColumnsMeta(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.Column, error)
	GetTableIndices(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.IndexMeta, error)
	GetTableReferences(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.ReferenceMeta, error)
}

type OffDbHandler interface {
	GetTablesList(ctx context.Context, conf *utils.ConfigYaml) ([]string, error)
	GetTableColumnsMeta(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.Column, error)
	GetTableIndices(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.IndexMeta, error)
	GetTableReferences(ctx context.Context, conf *utils.ConfigYaml, name string) ([]sch.ReferenceMeta, error)

	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(col_type string, val string) string

	CreateTableStatement(table *sch.TableMeta) string
	DropTableStatement(table *sch.TableMeta) string
	CreateColumnStatement(col *sch.Column) string
	AlterColumnStatement(col *sch.Column) string
	DropColumnStatement(col *sch.Column) string
	CreateIndexStatement(idx *sch.IndexMeta) string
	DropIndexStatement(idx *sch.IndexMeta) string
	CreateConstraintStatement(ref *sch.ReferenceMeta) string
	DropConstraintStatement(ref *sch.ReferenceMeta) string
}

func GetDriver(db_type string) (sch.DbHandler, error) {
	switch db_type {
	case "postgres":
		return &pg.PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

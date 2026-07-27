package schema

import (
	"context"

	utils "github.com/fastbear1/quack/internal"
)

type Meta interface {
	GetName() string
}

type AlterState struct {
	Type          uint8 // what shoul be altered 0 - data type, 1 - nullable, 2 - default value
	DataType      string
	IsNullable    bool
	ColumnDefault string
}

type Column struct {
	TableName         string
	ColumnName        string
	DataType          string
	IsNullable        bool
	ColumnDefault     string
	IsPrimary         bool
	PrimaryConstraint string
	AlterState        AlterState
}

type ReferenceMeta struct {
	TableName  string
	Name       string
	Column     string
	RefTable   string
	RefColumn  string
	RefOptions string
}

type IndexOption struct {
	Field      string
	Expression string
	Sort       string
	Collate    string
	Priority   int
}

type IndexMeta struct {
	TableName string
	Name      string
	Unique    bool
	Parsed    bool
	Columns   []IndexOption
	Type      string
	Where     string
	Option    string
}

type TableMeta struct {
	Name       string
	Columns    []Column
	References []ReferenceMeta
	Indeces    []IndexMeta
}

// Implement Meta interface
func (col Column) GetName() string {
	return col.ColumnName
}

func (ref ReferenceMeta) GetName() string {
	return ref.Name
}

func (idx IndexMeta) GetName() string {
	return idx.Name
}

func (t TableMeta) GetName() string {
	return t.Name
}

// Database driver interface
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

// SQL commands
func (table *TableMeta) CreateTable(drv DbHandler) string {
	return drv.CreateTableStatement(table)
}

func (table *TableMeta) DeleteTable(drv DbHandler) string {
	return drv.DropTableStatement(table)
}

func (col *Column) CreateColumn(drv DbHandler) string {
	return drv.CreateColumnStatement(col)
}

func (col *Column) AlterColumn(drv DbHandler) string {
	return drv.AlterColumnStatement(col)
}

func (col *Column) DeleteColumn(drv DbHandler) string {
	return drv.DropColumnStatement(col)
}

func (idx *IndexMeta) CreateIndex(drv DbHandler) string {
	return drv.CreateIndexStatement(idx)
}

func (idx *IndexMeta) DropIndex(drv DbHandler) string {
	return drv.DropIndexStatement(idx)
}

func (ref *ReferenceMeta) CreateConstraint(drv DbHandler) string {
	return drv.CreateConstraintStatement(ref)
}

func (ref *ReferenceMeta) DeleteConstraint(drv DbHandler) string {
	return drv.DropConstraintStatement(ref)
}

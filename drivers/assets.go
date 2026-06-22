package drivers

import (
	utils "github.com/fastbear1/quack/internal"
)

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
	Name         string
	Columns      []string
	IndexType    string
	IndexClass   string
	IndexWhere   string
	IndexExpr    string
	IndexSort    string
	IndexCollate string
	IndexOption  string
}

type TableMeta struct {
	Name       string
	Columns    []Column
	References []ReferenceMeta
	Indeces    []IndexMeta
}

func (table *TableMeta) CreateTable(conf *utils.ConfigYaml, drv DbHandler) (string, string) {
	return drv.CreateTableStatement(conf, table)
}

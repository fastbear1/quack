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
	Name    string
	Unique  bool
	Parsed  bool
	Columns []IndexOption
	Type    string
	Where   string
	Option  string
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

func (idx *IndexMeta) CreateIndex(conf *utils.ConfigYaml, drv DbHandler) (string, string) {
	return drv.CreateIndexStatement(conf, idx)
}

func (ref *ReferenceMeta) CreateConstraint(conf *utils.ConfigYaml, drv DbHandler) (string, string) {
	return drv.CreateConstraintStatement(conf, ref)
}

package schema

type Meta interface {
	GetName() string
}

type AlterState struct {
	Type          uint8 // what shoul be altered 0 - data type, 1 - nullable, 2 - default value
	DataType      string
	IsNullable    bool
	ColumnDefault string
	Processed     bool
}

type PrimaryOptions struct {
	IsSerial   bool
	IsIdentity bool
}

type Column struct {
	TableName      string
	ColumnName     string
	DataType       string
	IsNullable     bool
	ColumnDefault  string
	IsPrimary      bool
	PrimaryOptions PrimaryOptions
	AlterState     []AlterState
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
type DbInterface interface {
	GetTablesList() ([]string, error)
	GetTableColumnsMeta(name string) ([]Column, error)
	GetTableIndices(name string) ([]IndexMeta, error)
	GetTableReferences(name string) ([]ReferenceMeta, error)

	TransformName(name string) string
	TransformNull(nullable bool, def_val string) bool
	TransformType(g_type string) string
	TransformDefault(col_type string, val string) string
	TransformConstraintAction(action string) string

	CreateIndexName(table string, columns []string, exp string) string
	CreateConstraintName(table string, column string, refTable string, refColumn string) string

	CreateTableStatement(table *TableMeta) string
	DropTableStatement(table *TableMeta) string
	CreateColumnStatement(col *Column) string
	AlterColumnStatement(col *Column) string
	AlterDowngadeColumnStatement(col *Column) string
	DropColumnStatement(col *Column) string
	CreateIndexStatement(idx *IndexMeta) string
	DropIndexStatement(idx *IndexMeta) string
	CreateConstraintStatement(ref *ReferenceMeta) string
	DropConstraintStatement(ref *ReferenceMeta) string
}

// SQL commands
func (table *TableMeta) CreateTable(drv DbInterface) string {
	return drv.CreateTableStatement(table)
}

func (table *TableMeta) DeleteTable(drv DbInterface) string {
	return drv.DropTableStatement(table)
}

func (col *Column) CreateColumn(drv DbInterface) string {
	return drv.CreateColumnStatement(col)
}

func (col *Column) AlterColumn(drv DbInterface) string {
	return drv.AlterColumnStatement(col)
}

func (col *Column) AlterDowngradeColumn(drv DbInterface) string {
	return drv.AlterDowngadeColumnStatement(col)
}

func (col *Column) DeleteColumn(drv DbInterface) string {
	return drv.DropColumnStatement(col)
}

func (idx *IndexMeta) CreateIndex(drv DbInterface) string {
	return drv.CreateIndexStatement(idx)
}

func (idx *IndexMeta) DropIndex(drv DbInterface) string {
	return drv.DropIndexStatement(idx)
}

func (ref *ReferenceMeta) CreateConstraint(drv DbInterface) string {
	return drv.CreateConstraintStatement(ref)
}

func (ref *ReferenceMeta) DeleteConstraint(drv DbInterface) string {
	return drv.DropConstraintStatement(ref)
}

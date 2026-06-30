package drivers

type Column struct {
	TableName         string
	ColumnName        string
	DataType          string
	IsNullable        bool
	ColumnDefault     string
	IsPrimary         bool
	PrimaryConstraint string
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

func (table *TableMeta) CreateTable(drv DbHandler) (string, string) {
	return drv.CreateTableStatement(table)
}

func (col *Column) CreateColumn(drv DbHandler) (string, string) {
	return drv.CreateColumnStatement(col)
}

func (col *Column) DeleteColumn(drv DbHandler) (string, string) {
	return drv.DropColumnStatement(col)
}

func (idx *IndexMeta) CreateIndex(drv DbHandler) (string, string) {
	return drv.CreateIndexStatement(idx)
}

func (idx *IndexMeta) DropIndex(drv DbHandler) (string, string) {
	return drv.DropIndexStatement(idx)
}

func (ref *ReferenceMeta) CreateConstraint(drv DbHandler) (string, string) {
	return drv.CreateConstraintStatement(ref)
}

func (ref *ReferenceMeta) DeleteConstraint(drv DbHandler) (string, string) {
	return drv.DropConstraintStatement(ref)
}

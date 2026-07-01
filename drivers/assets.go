package drivers

type Meta interface {
	GetName() string
}

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

// Correspondance to Meta interface
func (col Column) GetName() string {
	return col.ColumnName
}

func (ref *ReferenceMeta) GetName() string {
	return ref.Name
}

func (idx *IndexMeta) GetName() string {
	return idx.Name
}

func (t *TableMeta) GetName() string {
	return t.Name
}

// SQL commands
func (table *TableMeta) CreateTable(drv DbHandler) (string, string) {
	return drv.CreateTableStatement(table)
}

func (table *TableMeta) DeleteTable(drv DbHandler) (string, string) {
	return drv.DropTableStatement(table)
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

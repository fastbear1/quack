package pgdriver

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"text/template"

	utils "github.com/fastbear1/quack/internal"
	. "github.com/fastbear1/quack/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Type conversion from Go type to postgres types
var TypeConversion = map[string]string{
	"int":     "bigint",
	"int16":   "smallint",
	"int32":   "bigint",
	"uint":    "bigint",
	"uint16":  "smallint",
	"uint32":  "bigint",
	"string":  "text",
	"float32": "real",
	"float64": "double precision",
}

// Templates function for calculate last item
var funcMap = template.FuncMap{
	"isLast": func(index int, len int) bool {
		return index+1 == len
	},
}

// Define connection variables
var (
	conn    *pgx.Conn
	connErr error
	once    sync.Once
)

// Database driver for postgres
type PgDriver struct {
	ctx  context.Context
	conf *utils.ConfigYaml
}

func GetPgDriver(ctx context.Context, conf *utils.ConfigYaml) *PgDriver {
	return &PgDriver{
		ctx:  ctx,
		conf: conf,
	}
}

var _ DbInterface = (*PgDriver)(nil)

// Get single connection
func getConnection(ctx context.Context, uri string) (*pgx.Conn, error) {
	dbCtx := ctx
	dbURI := uri
	once.Do(func() {
		conn, connErr = pgx.Connect(dbCtx, dbURI)
	})
	return conn, connErr
}

func (pg *PgDriver) GetTablesList() ([]string, error) {
	conn, err := getConnection(pg.ctx, pg.conf.Database.Uri.String())
	if err != nil {
		return []string{}, err
	}
	defer conn.Close(pg.ctx)
	// get tables list
	var tables []string
	rows, err := conn.Query(
		pg.ctx,
		GetTableNamesQuery,
		pgx.NamedArgs{
			"db": pg.conf.Database.Name,
		},
	)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Query error when getting tables list. %s", err)
		return []string{}, err
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		if utils.InArray(pg.conf.Database.Exclude, name) {
			fmt.Printf("Skipping table %s\n", name)
		} else {
			fmt.Printf("Found table %s\n", name)
			tables = append(tables, name)
		}
	}
	return tables, err
}

func (pg *PgDriver) GetTableColumnsMeta(name string) ([]Column, error) {
	var res = []Column{}
	conn, err := pgx.Connect(pg.ctx, string(pg.conf.Database.Uri))
	if err != nil {
		return []Column{}, err
	}
	defer conn.Close(pg.ctx)

	rows, err := conn.Query(
		pg.ctx,
		GetTableColumnsQuery,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer rows.Close()
	if err != nil {
		fmt.Println("Quering table columns metadata error...")
		return []Column{}, err
	}
	notes, err := pgx.CollectRows(rows, pgx.RowToStructByName[PgColumn])
	if err != nil {
		return []Column{}, err
	}

	pk_const_name, pk_column_name := pg.GetPrimaryKeyColumn(conn, pg.ctx, name)

	for i := 0; i < len(notes); i++ {
		res = append(res, Column{
			TableName:     name,
			ColumnName:    notes[i].Column_name,
			DataType:      normalizeCharacterVariyng(notes[i].Data_type, notes[i].Character_maximum_length),
			IsNullable:    transformNullToString(notes[i].Is_nullable),
			ColumnDefault: notes[i].Column_default.String,
			IsPrimary: func(lname string, rname string) bool {
				if lname == rname {
					return true
				}
				return false
			}(notes[i].Column_name, pk_column_name),
			PrimaryConstraint: func(lname string, rname string) string {
				if lname == rname {
					return pk_const_name
				}
				return ""
			}(notes[i].Column_name, pk_column_name),
		})
	}
	return res, nil
}

func (pg *PgDriver) GetPrimaryKeyColumn(conn *pgx.Conn, ctx context.Context, table_name string) (string, string) {
	// Find primary key field
	rowid, err := conn.Query(
		ctx,
		FindPrimaryKeyQuery,
		pgx.NamedArgs{
			"table": table_name,
		},
	)
	defer rowid.Close()
	utils.CheckErrLite(err)

	var pk_const_name, pk_column_name string
	for rowid.Next() {
		err = rowid.Scan(&pk_const_name, &pk_column_name)
		utils.CheckErrLite(err)
	}

	return pk_const_name, pk_column_name
}

func (pg *PgDriver) GetTableIndices(name string) ([]IndexMeta, error) {
	conn, err := pgx.Connect(pg.ctx, string(pg.conf.Database.Uri))
	var idxt []IndexMeta

	if err != nil {
		fmt.Println("Quering indices definition data error...")
		return []IndexMeta{}, err
	}
	defer conn.Close(pg.ctx)

	row, err := conn.Query(
		pg.ctx,
		GetTableIndicesInformation,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer row.Close()
	if err != nil {
		fmt.Println("Quering table indices error...")
		return []IndexMeta{}, err
	}

	var idx_const_name, idx_const_def string
	for row.Next() {
		err = row.Scan(&idx_const_name, &idx_const_def)
		utils.CheckErrLite(err)
		idx, _ := ParseDatabaseIndices(idx_const_def)
		idxt = append(idxt, idx)
	}

	return idxt, nil
}

func (pg *PgDriver) GetTableReferences(name string) ([]ReferenceMeta, error) {
	conn, err := pgx.Connect(pg.ctx, string(pg.conf.Database.Uri))
	var ref []ReferenceMeta

	if err != nil {
		fmt.Println("Quering constraint definition data error...")
		return []ReferenceMeta{}, err
	}
	defer conn.Close(pg.ctx)

	row, err := conn.Query(
		pg.ctx,
		GetTableForeignKeys,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer row.Close()
	if err != nil {
		fmt.Println("Quering table constraints error...")
		return []ReferenceMeta{}, err
	}
	var ref_name, ref_const_def string
	for row.Next() {
		err = row.Scan(&ref_name, &ref_const_def)
		utils.CheckErrLite(err)
		res, _ := ParseDatabaseReferences(name, ref_name, ref_const_def)
		ref = append(ref, res)
	}
	return ref, nil
}

func (pg *PgDriver) TransformName(name string) string {
	// Camel case to snake case
	var buffer bytes.Buffer
	delta := 'a' - 'A'
	for i, v := range name {
		if i == 0 && rune(v) < 'a' {
			buffer.WriteRune(rune(v + delta))
		} else if i > 0 && rune(name[i-1]) >= 'a' && rune(v) < 'a' {
			buffer.WriteRune(rune('_'))
			buffer.WriteRune(rune(v + delta))
		} else {
			if rune(v) < 'a' {
				v = v + delta
			}
			buffer.WriteRune(rune(v))
		}
	}
	return buffer.String()
}

func (pg *PgDriver) TransformNull(nullable bool, def_val string) bool {
	var use_null bool = false
	if def_val == "" && !nullable {
		use_null = true
	}
	return use_null
}

func (pg *PgDriver) TransformType(g_type string) string {
	tp, ok := TypeConversion[g_type]
	if !ok {
		return g_type
	}
	return tp
}

func (pg *PgDriver) TransformDefault(columnType string, columnDefault string) string {
	defValue := columnDefault
	r := regexp.MustCompile(`varchar\(\d+\)`)
	if r.MatchString(columnType) {
		defValue = fmt.Sprintf("'%s'::text", defValue)
	}

	switch columnType {
	case "text":
		defValue = fmt.Sprintf("'%s'::text", defValue)
	}
	return defValue
}

func (pg *PgDriver) CreateIndexName(table string, columns []string, exp string) string {
	indexName := "idx_"
	columnsSuffix := ""
	expSuffix := ""
	if exp != "" {
		expParts := strings.Split(exp, "(")
		expSuffix += fmt.Sprintf("_%s", expParts[0])
	}
	for _, c := range columns {
		columnsSuffix += fmt.Sprintf("_%s", c)
	}
	indexName += table + columnsSuffix + expSuffix
	return indexName
}

func (pg *PgDriver) CreateConstraintName(table string, column string, refTable string, refColumn string) string {
	return fmt.Sprintf("fk_%s_%s__%s_%s", table, column, refTable, refColumn)
}

func (pg *PgDriver) TransformConstraintAction(action string) string {
	defaultAction := "ON DELETE"
	if action == "OnUpdate" {
		defaultAction = "ON UPDATE"
	}
	return defaultAction
}

func (pg *PgDriver) CreateTableStatement(t *TableMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateTableTmpl)
	utils.CheckErrLite(err)

	// find primary column
	primary := ""
	for _, c := range t.Columns {
		if c.IsPrimary {
			primary = c.ColumnName
		}
	}
	var ft = struct {
		PrimaryColumn string
		*TableMeta
	}{
		primary,
		t,
	}

	if err := masterTmpl.Execute(&sqlCommand, ft); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) DropTableStatement(t *TableMeta) string {
	var sqlCommand bytes.Buffer

	deleteTmpl, err := template.New("delete").Parse(DropTableTmpl)
	utils.CheckErrLite(err)

	if err := deleteTmpl.Execute(&sqlCommand, t); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) CreateColumnStatement(col *Column) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateColumnTmpl)
	utils.CheckErrLite(err)

	if err := masterTmpl.Execute(&sqlCommand, col); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) AlterColumnStatement(col *Column) string {
	return getAlterColumnCommand(col, false)
}

func (pg *PgDriver) DropColumnStatement(col *Column) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropColumnTmpl)
	utils.CheckErrLite(err)

	if err := masterTmpl.Execute(&sqlCommand, col); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func getAlterColumnCommand(col any, downgrade bool) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(AlterColumnTmpl)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, col); err != nil {
		fmt.Println(err)
	}
	sql := sqlCommand.String()

	if downgrade {
		alt := col.(*AlterData)
		switch alt.Type {
		case 0:
			sql = sql + " " + fmt.Sprintf("TYPE %s", alt.DataType)
		case 1:
			if alt.IsNullable == true {
				sql = sql + " " + "DROP NOT NULL"
			} else {
				sql = sql + " " + "SET NOT NULL"
			}
		case 2:
			if alt.ColumnDefault == "" {
				sql = sql + " " + "DROP DEFAULT"
			} else {
				sql = sql + " " + fmt.Sprintf("SET DEFAULT %s", alt.ColumnDefault)
			}
		}
	} else {
		alt := col.(*Column)
		switch alt.AlterState.Type {
		case 0:
			sql = sql + " " + fmt.Sprintf("TYPE %s", alt.DataType)
		case 1:
			if alt.IsNullable == true {
				sql = sql + " " + "DROP NOT NULL"
			} else {
				sql = sql + " " + "SET NOT NULL"
			}
		case 2:
			if alt.ColumnDefault == "" {
				sql = sql + " " + "DROP DEFAULT"
			} else {
				sql = sql + " " + fmt.Sprintf("SET DEFAULT %s", alt.ColumnDefault)
			}
		}
	}
	return sql
}

func (pg *PgDriver) CreateIndexStatement(idx *IndexMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateIndexTmpl)
	utils.CheckErrLite(err)

	// prapre columns names for composite index
	var fields string
	flen := len(idx.Columns)

	switch {
	case flen > 1:
		var farray []string
		utils.SortArray(idx.Columns, func(i, j int) bool {
			return idx.Columns[i].Priority > idx.Columns[j].Priority
		})
		for _, c := range idx.Columns {
			farray = append(farray, c.Field)
		}
		fields = strings.Join(farray, ", ")
	default:
		fields = idx.Columns[0].Field
	}

	// prepare expression if exists
	var expStr string
	if idx.Columns[0].Expression != "" {
		expStr = strings.Split(idx.Columns[0].Expression, "(")[0]
	}

	var t = struct {
		TableName  string
		Name       string
		Unique     bool
		Type       string
		Expression string
		Columns    string
	}{
		idx.TableName,
		idx.Name,
		idx.Unique,
		idx.Type,
		expStr,
		fields,
	}

	if err := masterTmpl.Execute(&sqlCommand, t); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) DropIndexStatement(idx *IndexMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropIndexTmpl)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, idx); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) CreateConstraintStatement(ref *ReferenceMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateConstraintTmpl)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, ref); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgDriver) DropConstraintStatement(ref *ReferenceMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropConstraintTmpl)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, ref); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

// helpers
func normalizeCharacterVariyng(data_type string, lenght pgtype.Uint32) string {
	if data_type == "character varying" {
		data_type = fmt.Sprintf("varchar(%d)", lenght.Uint32)
	}
	return data_type
}

func transformNullToString(isnull string) bool {
	if isnull == "NO" {
		return false
	}
	return true
}

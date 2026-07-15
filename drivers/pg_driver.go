package drivers

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	utils "github.com/fastbear1/quack/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Internla Pg type definitions
type PgColumn struct {
	Column_name              string
	Data_type                string
	Character_maximum_length pgtype.Uint32
	Is_nullable              string
	Column_default           pgtype.Text
}

type AlterData struct {
	TableName     string
	ColumnName    string
	Type          uint8
	DataType      string
	IsNullable    bool
	ColumnDefault string
}

// Type conversion from Go type to postgres types
var TypeConversion = map[string]string{
	"uint":    "bigint",
	"uint16":  "smallint",
	"string":  "text",
	"float32": "real",
	"float64": "double precision",
}

const (
	GetTableColumnsQuery = `
SELECT 
	column_name, 
	data_type, 
	character_maximum_length, 
	is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = @table`
	FindPrimaryKeyQuery = `
SELECT 
	kc.constraint_name, 
	kc.column_name 
FROM information_schema.key_column_usage kc 
JOIN information_schema.table_constraints tc 
	ON kc.constraint_name = tc.constraint_name 
WHERE tc.constraint_type = 'PRIMARY KEY' 
	AND kc.table_name=@table`
	GetTableIndicesInformation = `
SELECT 
	indexname,
	indexdef 
FROM pg_catalog.pg_indexes 
WHERE tablename=@table 
	AND indexname NOT IN (
		SELECT constraint_name 
		FROM information_schema.table_constraints 
		WHERE table_name=@table
);`
	GetTableNamesQuery = `
SELECT 
	table_name 
FROM information_schema.tables 
WHERE table_type='BASE TABLE' 
	AND table_schema='public' 
	AND table_catalog=@db`
	GetTableForeignKeys = `
SELECT 
	conname, 
	pg_get_constraintdef(oid) 
FROM pg_constraint
WHERE contype IN ('f', 'p ')
	AND pg_get_constraintdef(oid) LIKE 'FOREIGN KEY %' 
	AND conrelid::regclass::text = @table;
`
)

// SQL templates and functions
const (
	CreateTemaplete = `{{$lenColumns := len .Columns}}{{$lenRef := len .References}}CREATE TABLE "public"."{{ .Name }}"(
{{- range $i, $a := .Columns}}
	{{ .ColumnName }} {{ .DataType }}{{if not .IsNullable}} NOT NULL{{end}}{{ if .ColumnDefault }} default {{ .ColumnDefault }}{{ end }},
{{- end}}
	{{ if .PrimaryColumn }}PRIMARY KEY ("{{.PrimaryColumn}}"){{ end }}{{ if .References }},{{end}}
{{- range $i, $a := .References}}
	CONSTRAINT "{{.Name}}" FOREIGN KEY ("{{.Column}}") REFERENCES "public"."{{.RefTable}}" ("{{.RefColumn}}"){{if .RefOptions}} {{.RefOptions}}{{end}}{{ if not (isLast $i $lenRef) }},{{ end }}
{{- end}}
);`
	DropTableTemplate = `DROP TABLE IF EXISTS "public"."{{.Name}}";`
	CreateColumn      = `ALTER TABLE "public"."{{.TableName}}" ADD COLUMN IF NOT EXISTS {{ .ColumnName }} {{ .DataType }}{{if not .IsNullable}} NOT NULL{{end}}{{ if .ColumnDefault }} default {{ .ColumnDefault }}{{ end }}`
	AlterColumn       = `ALTER TABLE "public"."{{.TableName}}" ALTER COLUMN IF EXISTS {{ .ColumnName }}`
	DropColumn        = `ALTER TABLE "public"."{{.TableName}}" DROP COLUMN IF EXISTS {{ .ColumnName }}`
	CreateIndex       = `CREATE INDEX IF NOT EXISTS "{{.Name}}" ON "public"."{{.TableName}}"{{if .Unique}} UNIQUE{{end}} USING {{.Type}} {{.Expression}}({{.Columns}});`
	DropIndex         = `DROP INDEX IF EXISTS "{{.Name}}"`
	CreateConstraint  = `ALTER TABLE "public"."{{.TableName}}" ADD CONSTRAINT IF NOT EXISTS "{{.Name}}" FOREIGN KEY ("{{.Column}}") REFERENCES "public"."{{.RefTable}}" ("{{.RefColumn}}"){{if .RefOptions}} {{.RefOptions}}{{end}}`
	DropConstraint    = `ALTER TABLE "public"."{{.TableName}}" DROP CONSTRAINT IF EXISTS "{{.Name}}"`
)

var funcMap = template.FuncMap{
	"isLast": func(index int, len int) bool {
		return index+1 == len
	},
}

// Database driver for postgres
type PgHandler struct{}

func (pg *PgHandler) GetTablesList(conf *utils.ConfigYaml) ([]string, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, string(conf.Database.Uri))
	if err != nil {
		return []string{}, err
	}
	defer conn.Close(ctx)

	// get tables list
	dbtables, err := getDbTables(conf, ctx, conn)
	if err != nil {
		return []string{}, nil
	}
	return dbtables, nil
}

func getDbTables(conf *utils.ConfigYaml, ctx context.Context, conn *pgx.Conn) ([]string, error) {
	var tables []string
	rows, err := conn.Query(
		ctx,
		GetTableNamesQuery,
		pgx.NamedArgs{
			"db": conf.Database.Name,
		},
	)
	defer rows.Close()
	if err != nil {
		fmt.Printf("Query error when getting tables list. %s", err)
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		if utils.InArray(conf.Database.Exclude, name) {
			fmt.Printf("Skipping table %s\n", name)
		} else {
			fmt.Printf("Found table %s\n", name)
			tables = append(tables, name)
		}
	}
	return tables, err
}

func (pg *PgHandler) GetTableColumnsMeta(conf *utils.ConfigYaml, name string) ([]Column, error) {
	var res = []Column{}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, string(conf.Database.Uri))
	if err != nil {
		return []Column{}, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(
		ctx,
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

	pk_const_name, pk_column_name := pg.GetPrimaryKeyColumn(conn, ctx, name)

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

func (pg *PgHandler) GetPrimaryKeyColumn(conn *pgx.Conn, ctx context.Context, table_name string) (string, string) {
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

func (pg *PgHandler) GetTableIndices(conf *utils.ConfigYaml, name string) ([]IndexMeta, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, string(conf.Database.Uri))
	var idxt []IndexMeta

	if err != nil {
		return []IndexMeta{}, err
	}
	defer conn.Close(ctx)

	row, err := conn.Query(
		ctx,
		GetTableIndicesInformation,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer row.Close()
	utils.CheckErrLite(err)

	var idx_const_name, idx_const_def string
	for row.Next() {
		err = row.Scan(&idx_const_name, &idx_const_def)
		utils.CheckErrLite(err)
		idx, err := ParseDatabaseIndices(idx_const_def)
		utils.CheckErrLite(err)
		idxt = append(idxt, idx)
	}

	return idxt, nil
}

func (pg *PgHandler) GetTableReferences(conf *utils.ConfigYaml, name string) ([]ReferenceMeta, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, string(conf.Database.Uri))
	var ref []ReferenceMeta

	if err != nil {
		return []ReferenceMeta{}, err
	}
	defer conn.Close(ctx)

	row, err := conn.Query(
		ctx,
		GetTableForeignKeys,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer row.Close()
	utils.CheckErrLite(err)

	var ref_name, ref_const_def string
	for row.Next() {
		err = row.Scan(&ref_name, &ref_const_def)
		utils.CheckErrLite(err)
		res, err := ParseDatabaseReferences(name, ref_name, ref_const_def)
		utils.CheckErrLite(err)
		ref = append(ref, res)
	}

	return ref, nil
}

func (pg *PgHandler) TransformName(name string) string {
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

func (pg *PgHandler) TransformNull(nullable bool, def_val string) bool {
	var use_null bool = false
	if def_val == "" && !nullable {
		use_null = true
	}
	return use_null
}

func (pg *PgHandler) TransformType(g_type string) string {
	tp, ok := TypeConversion[g_type]
	if !ok {
		return g_type
	}
	return tp
}

func (pg *PgHandler) TransformDefault(columnType string, columnDefault string) string {
	defValue := columnDefault
	columnType = strings.Split(columnType, "(")[0]
	switch columnType {
	case "varchar":
		defValue = fmt.Sprintf("'%s'::text", defValue)
	}
	return defValue
}

func (pg *PgHandler) CreateTableStatement(t *TableMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateTemaplete)
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

func (pg *PgHandler) DropTableStatement(t *TableMeta) string {
	var sqlCommand bytes.Buffer

	deleteTmpl, err := template.New("delete").Parse(DropTableTemplate)
	utils.CheckErrLite(err)

	if err := deleteTmpl.Execute(&sqlCommand, t); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgHandler) CreateColumnStatement(col *Column) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateColumn)
	utils.CheckErrLite(err)

	if err := masterTmpl.Execute(&sqlCommand, col); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgHandler) AlterColumnStatement(col *Column) string {
	return getAlterColumnCommand(col, false)

	// declate temporary strict for downgrade alter command

	/*
		var data = AlterData{
			TableName:     col.TableName,
			ColumnName:    col.ColumnName,
			Type:          col.AlterState.Type,
			DataType:      col.AlterState.DataType,
			IsNullable:    col.AlterState.IsNullable,
			ColumnDefault: col.AlterState.ColumnDefault,
		}
		sqlDown = getAlterColumnCommand(&data, true)
	*/
}

func (pg *PgHandler) DropColumnStatement(col *Column) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropColumn)
	utils.CheckErrLite(err)

	if err := masterTmpl.Execute(&sqlCommand, col); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func getAlterColumnCommand(col any, downgrade bool) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(AlterColumn)
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

func (pg *PgHandler) CreateIndexStatement(idx *IndexMeta) string {
	//TODO: to refactor
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateIndex)
	utils.CheckErrLite(err)

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
		idx.Columns[0].Expression,
		idx.Columns[0].Field,
	}

	if err := masterTmpl.Execute(&sqlCommand, t); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgHandler) DropIndexStatement(idx *IndexMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropIndex)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, idx); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgHandler) CreateConstraintStatement(ref *ReferenceMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(CreateConstraint)
	utils.CheckErrLite(err)
	if err := masterTmpl.Execute(&sqlCommand, ref); err != nil {
		fmt.Println(err)
	}
	return sqlCommand.String()
}

func (pg *PgHandler) DropConstraintStatement(ref *ReferenceMeta) string {
	var sqlCommand bytes.Buffer
	masterTmpl, err := template.New("master").Funcs(funcMap).Parse(DropConstraint)
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

package drivers

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	utils "github.com/fastbear1/quack/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Type definitions
type PgColumn struct {
	Column_name              string
	Data_type                string
	Character_maximum_length pgtype.Uint32
	Is_nullable              string
	Column_default           pgtype.Text
}

// Type conversion
var TypeConversion = map[string]string{
	"uint":   "bigint",
	"uint16": "smallint",
}

const (
	CreateTemaplete = `CREATE TABLE "public"."{{ .TableName }}"(
{{- range .Columns}}
   {{.}}
{{- end}}
);`
	DropTableTemplate = `DROP TABLE IF EXISTS "public"."{{.TableName}}";`
	CreateColumn      = `{{ .ColName }} {{ .ColType }}{{ if .ColPrimary }} PRIMARY KEY{{ end }}{{if .UseNull}} NOT NULL{{end}}{{ if .ColDefault }} default {{ .ColDefault }}{{ end }}{{ if not .LastColumn }},{{ end }}`
)

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

func (pg *PgHandler) GetTableColumnsMeta(conf *utils.ConfigYaml, name string) ([]Column, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, string(conf.Database.Uri))
	var res []Column = []Column{}
	if err != nil {
		return []Column{}, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(
		ctx,
		`SELECT column_name, data_type, character_maximum_length, is_nullable, column_default 
		FROM information_schema.columns 
		WHERE table_name = @table`,
		pgx.NamedArgs{
			"table": name,
		},
	)
	if err != nil {
		fmt.Println("Quering table columns metadata error...")
		return []Column{}, err
	}
	notes, err := pgx.CollectRows(rows, pgx.RowToStructByName[PgColumn])
	if err != nil {
		return []Column{}, err
	}
	// Find primary key field
	rowid, err := conn.Query(
		ctx,
		`SELECT 
			kc.constraint_name, 
			kc.column_name 
		FROM information_schema.key_column_usage kc 
		JOIN information_schema.table_constraints tc 
			ON kc.constraint_name = tc.constraint_name 
		WHERE tc.constraint_type = 'PRIMARY KEY' 
			AND kc.table_name=@table`,
		pgx.NamedArgs{
			"table": name,
		},
	)
	defer rowid.Close()

	utils.CheckErrLite(err)
	var pk_const_name, pk_column_name string

	for rowid.Next() {
		err = rowid.Scan(&pk_const_name, &pk_column_name)
		utils.CheckErrLite(err)
	}

	for i := 0; i < len(notes); i++ {
		res = append(res, Column{
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

func normalizeCharacterVariyng(data_type string, lenght pgtype.Uint32) string {
	if data_type == "character varying" {
		data_type = fmt.Sprintf("varchar(%d)", lenght.Uint32)
	}
	return data_type
}

func getDbTables(conf *utils.ConfigYaml, ctx context.Context, conn *pgx.Conn) ([]string, error) {
	var tables []string
	rows, err := conn.Query(
		ctx,
		"SELECT table_name FROM information_schema.tables WHERE table_type='BASE TABLE' AND table_schema='public' AND table_catalog=@db",
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

func transformNullToString(isnull string) bool {
	if isnull == "NO" {
		return false
	}
	return true
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

func (pg *PgHandler) TransformDefault(val string) string {
	defval := val
	if defval != "" {
		match, _ := regexp.MatchString(`(.*)\(\)(.*)`, val)
		if !match {
			// should inherit column type
			defval = fmt.Sprintf("'%s'", val)
		}
	}
	return defval
}

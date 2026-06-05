package drivers

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	utils "github.com/fastbear1/quack/internal"
	"github.com/jackc/pgx/v5"
)

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

func (pg *PgHandler) GetData(conf *utils.ConfigYaml) ([]string, error) {
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

func TransformName(name string) string {
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

func TransformNull(nullable bool, def_val string) bool {
	var use_null bool = false
	if def_val == "" && !nullable {
		use_null = true
	}
	return use_null
}

func TransformType(g_type string) string {
	tp, ok := TypeConversion[g_type]
	if !ok {
		return g_type
	}
	return tp
}

func TransformDefault(val string) string {
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

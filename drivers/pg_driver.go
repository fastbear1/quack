package drivers

import (
	"bytes"
	"fmt"
	"regexp"
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

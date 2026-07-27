package pgdriver

// query templates
const (
	GetTableNamesQuery = `
SELECT 
	table_name 
FROM information_schema.tables 
WHERE table_type='BASE TABLE' 
	AND table_schema='public' 
	AND table_catalog=@db`
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
	GetTableForeignKeys = `
SELECT 
	conname, 
	pg_get_constraintdef(oid) 
FROM pg_constraint
WHERE contype IN ('f', 'p ')
	AND pg_get_constraintdef(oid) LIKE 'FOREIGN KEY %' 
	AND conrelid::regclass::text = @table;`
)

// SQL commands templates
const (
	CreateTableTmpl = `{{$lenColumns := len .Columns}}{{$lenRef := len .References}}CREATE TABLE "public"."{{ .Name }}"(
{{- range $i, $a := .Columns}}
	{{ .ColumnName }} {{ .DataType }}{{if not .IsNullable}} NOT NULL{{end}}{{ if .ColumnDefault }} default {{ .ColumnDefault }}{{ end }},
{{- end}}
	{{ if .PrimaryColumn }}PRIMARY KEY ("{{.PrimaryColumn}}"){{ end }}{{ if .References }},{{end}}
{{- range $i, $a := .References}}
	CONSTRAINT "{{.Name}}" FOREIGN KEY ("{{.Column}}") REFERENCES "public"."{{.RefTable}}" ("{{.RefColumn}}"){{if .RefOptions}} {{.RefOptions}}{{end}}{{ if not (isLast $i $lenRef) }},{{ end }}
{{- end}}
);`
	DropTableTmpl        = `DROP TABLE IF EXISTS "public"."{{.Name}}";`
	CreateColumnTmpl     = `ALTER TABLE "public"."{{.TableName}}" ADD COLUMN IF NOT EXISTS {{ .ColumnName }} {{ .DataType }}{{if not .IsNullable}} NOT NULL{{end}}{{ if .ColumnDefault }} default {{ .ColumnDefault }}{{ end }}`
	AlterColumnTmpl      = `ALTER TABLE "public"."{{.TableName}}" ALTER COLUMN IF EXISTS {{ .ColumnName }}`
	DropColumnTmpl       = `ALTER TABLE "public"."{{.TableName}}" DROP COLUMN IF EXISTS {{ .ColumnName }}`
	CreateIndexTmpl      = `CREATE INDEX IF NOT EXISTS "{{.Name}}" ON "public"."{{.TableName}}"{{if .Unique}} UNIQUE{{end}} USING {{.Type}} {{.Expression}}({{.Columns}});`
	DropIndexTmpl        = `DROP INDEX IF EXISTS "{{.Name}}"`
	CreateConstraintTmpl = `ALTER TABLE "public"."{{.TableName}}" ADD CONSTRAINT IF NOT EXISTS "{{.Name}}" FOREIGN KEY ("{{.Column}}") REFERENCES "public"."{{.RefTable}}" ("{{.RefColumn}}"){{if .RefOptions}} {{.RefOptions}}{{end}}`
	DropConstraintTmpl   = `ALTER TABLE "public"."{{.TableName}}" DROP CONSTRAINT IF EXISTS "{{.Name}}"`
)

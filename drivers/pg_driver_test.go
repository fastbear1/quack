package drivers

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestTransformName(t *testing.T) {
	pg := &PgHandler{}

	var test = []struct {
		name   string
		expect string
	}{
		{"SimpleName", "simple_name"},
		{"NameWithCapitalAtEnD", "name_with_capital_at_en_d"},
		{"UPPERName", "uppername"},
		{"lowerToUp", "lower_to_up"},
	}
	for _, tt := range test {
		t.Run(fmt.Sprintf("Test for %s", tt.name), func(t *testing.T) {
			res := pg.TransformName(tt.name)
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestTransformNull(t *testing.T) {
	var test = []struct {
		name   string
		null   bool
		defval string
		expect bool
	}{
		{"Not null value", false, "not null", false},
		{"Use null value", false, "", true},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			res := (&PgHandler{}).TransformNull(tt.null, tt.defval)
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestTransformType(t *testing.T) {
	var test = []struct {
		name     string
		codeType string
		expect   string
	}{
		{"Uint type", "uint", "bigint"},
		{"Int64 type", "int64", "int64"},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			res := (&PgHandler{}).TransformType(tt.codeType)
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestNormalizeVarChar(t *testing.T) {
	var test = []struct {
		name     string
		datatype string
		lenght   uint32
		expect   string
	}{
		{"Default Varying Character", "character varying", 255, "varchar(255)"},
		{"Small Varying Character", "character varying", 10, "varchar(10)"},
		{"Not A Varying Character", "smallint", 100, "smallint"},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			res := normalizeCharacterVariyng(tt.datatype, pgtype.Uint32{Uint32: tt.lenght, Valid: true})
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestTransformNullToString(t *testing.T) {
	var test = []struct {
		name   string
		isnull string
		expect bool
	}{
		{"Is Null", "YES", true},
		{"Not NUll", "NO", false},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			res := transformNullToString(tt.isnull)
			assert.Equal(t, res, tt.expect)
		})
	}
}

// Tests for SQL creation methods
func TestCreaetTabelStatement(t *testing.T) {
	// &{simple_table [{id uuid false gen_random_uuid() true } {name varchar(255) false  false } {sid smallint false  false } {email varchar(255) false  false } {status varchar(10) false active false } {name_t varchar(255) false  false } {created_at timestamp false now() false } {updated_at timestamp false now() false }] [] []}
	var test = []struct {
		name      string
		tablemeta TableMeta
		expect    []string
	}{
		{
			name: "simple test for create table SQL",
			tablemeta: TableMeta{
				Name: "test_table",
				Columns: []Column{
					{
						ColumnName:    "id",
						DataType:      "uuid",
						IsNullable:    false,
						ColumnDefault: "gen_random_uuid()",
						IsPrimary:     true,
					},
					{
						ColumnName:    "name",
						DataType:      "varchar(255)",
						IsNullable:    false,
						ColumnDefault: "",
						IsPrimary:     false,
					},
					{
						ColumnName:    "status",
						DataType:      "varchar(10)",
						IsNullable:    false,
						ColumnDefault: "active",
						IsPrimary:     false,
					},
					{
						ColumnName:    "created_at",
						DataType:      "timestamp",
						IsNullable:    false,
						ColumnDefault: "now()",
						IsPrimary:     false,
					},
				},
				Indeces:    []IndexMeta{},
				References: []ReferenceMeta{},
			},
			expect: []string{
				`CREATE TABLE "public"."test_table"(
	id uuid NOT NULL default gen_random_uuid(),
	name varchar(255) NOT NULL,
	status varchar(10) NOT NULL default active,
	created_at timestamp NOT NULL default now(),
	PRIMARY KEY ("id")
);`,
				`DROP TABLE IF EXISTS "public"."test_table";`,
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			sqlUp, sqlDown := (&PgHandler{}).CreateTableStatement(&tt.tablemeta)
			assert.Equal(t, sqlUp, tt.expect[0])
			assert.Equal(t, sqlDown, tt.expect[1])
		})
	}
}

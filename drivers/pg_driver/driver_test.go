package pgdriver

import (
	"fmt"
	"testing"

	utils "github.com/fastbear1/quack/internal"
	. "github.com/fastbear1/quack/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestTransformName(t *testing.T) {
	pg := &PgDriver{}

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
			res := (&PgDriver{}).TransformNull(tt.null, tt.defval)
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
			res := (&PgDriver{}).TransformType(tt.codeType)
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestTransformDefault(t *testing.T) {
	var test = []struct {
		columnType    string
		columnDefault string
		expect        string
	}{
		{
			"varchar(10)",
			"active",
			"'active'::text",
		},
		{
			"text",
			"default text",
			"'default text'::text",
		},
	}

	for n, tt := range test {
		t.Run(fmt.Sprintf("Testing transformation default value #%d\n", n), func(t *testing.T) {
			res := (&PgDriver{}).TransformDefault(tt.columnType, tt.columnDefault)
			assert.Equal(t, res, tt.expect)
		})
	}
}

func TestNormalizeVarChar(t *testing.T) {
	var test = []struct {
		name     string
		datatype string
		lenght   uint32
		udt_name string
		expect   string
	}{
		{"Default Varying Character", "character varying", 255, "string", "varchar(255)"},
		{"Small Varying Character", "character varying", 10, "string", "varchar(10)"},
		{"Not A Varying Character", "smallint", 100, "int", "smallint"},
		{"Integer array type", "ARRAY", 10, "_int4", "integer[]"},
		{"Text array type", "ARRAY", 10, "_text", "text[]"},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			res := normalizeColumnDataType(tt.datatype, pgtype.Uint32{Uint32: tt.lenght, Valid: true}, tt.udt_name)
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

// Tests for SQL commands
func TestCreaetTabelStatement(t *testing.T) {
	var test = []struct {
		tablemeta TableMeta
		expect    []string
	}{
		{
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
		{
			tablemeta: TableMeta{
				Name: "simple_table",
				Columns: []Column{
					{
						TableName:      "simple_table",
						ColumnName:     "id",
						DataType:       "uuid",
						IsNullable:     false,
						ColumnDefault:  "gen_random_uuid()",
						IsPrimary:      true,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "name",
						DataType:       "varchar(255)",
						IsNullable:     false,
						ColumnDefault:  "",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "sid",
						DataType:       "smallint",
						IsNullable:     false,
						ColumnDefault:  "",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "email",
						DataType:       "varchar(255)",
						IsNullable:     false,
						ColumnDefault:  "",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "user_id",
						DataType:       "uuid",
						IsNullable:     false,
						ColumnDefault:  "",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "status",
						DataType:       "varchar(10)",
						IsNullable:     false,
						ColumnDefault:  "'active'::text",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "name_t",
						DataType:       "varchar(255)",
						IsNullable:     false,
						ColumnDefault:  "",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "created_at",
						DataType:       "timestamp",
						IsNullable:     false,
						ColumnDefault:  "now()",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
					{
						TableName:      "simple_table",
						ColumnName:     "updated_at",
						DataType:       "timestamp",
						IsNullable:     false,
						ColumnDefault:  "now()",
						IsPrimary:      false,
						PrimaryOptions: PrimaryOptions{},
					},
				},
				References: []ReferenceMeta{},
				Indeces:    []IndexMeta{},
			},
			expect: []string{
				`CREATE TABLE "public"."simple_table"(
	id uuid NOT NULL default gen_random_uuid(),
	name varchar(255) NOT NULL,
	sid smallint NOT NULL,
	email varchar(255) NOT NULL,
	user_id uuid NOT NULL,
	status varchar(10) NOT NULL default 'active'::text,
	name_t varchar(255) NOT NULL,
	created_at timestamp NOT NULL default now(),
	updated_at timestamp NOT NULL default now(),
	PRIMARY KEY ("id")
);`,
				`DROP TABLE IF EXISTS "public"."simple_table";`,
			},
		},
	}

	conf := utils.ConfigYaml{}
	conf.ReadConfig()

	for n, tt := range test {
		t.Run(fmt.Sprintf("ing method for creating table SQL command #%d", n), func(t *testing.T) {
			drv := PgDriver{}
			sqlUp := (&drv).CreateTableStatement(&tt.tablemeta)
			sqlDown := (&drv).DropTableStatement(&tt.tablemeta)
			assert.Equal(t, sqlUp, tt.expect[0])
			assert.Equal(t, sqlDown, tt.expect[1])
		})
	}
}

func TestColumnStatements(t *testing.T) {
	var testData = []struct {
		column   Column
		expected []string
	}{
		{
			Column{
				TableName:     "simple_table",
				ColumnName:    "sid",
				DataType:      "smallint",
				IsNullable:    false,
				ColumnDefault: "",
			},
			[]string{
				`ALTER TABLE "public"."simple_table" ADD COLUMN IF NOT EXISTS sid smallint NOT NULL;`,
				`ALTER TABLE "public"."simple_table" DROP COLUMN IF EXISTS sid;`,
			},
		},
		{

			Column{
				TableName:     "simple_table",
				ColumnName:    "email",
				DataType:      "varchar(255)",
				IsNullable:    false,
				ColumnDefault: "'test@mail.com'::text",
			},
			[]string{
				`ALTER TABLE "public"."simple_table" ADD COLUMN IF NOT EXISTS email varchar(255) NOT NULL DEFAULT 'test@mail.com'::text;`,
				`ALTER TABLE "public"."simple_table" DROP COLUMN IF EXISTS email;`,
			},
		},
		{
			Column{
				TableName:     "simple_table",
				ColumnName:    "user_id",
				DataType:      "uuid",
				IsNullable:    false,
				ColumnDefault: "gen_random_uuid()",
			},
			[]string{
				`ALTER TABLE "public"."simple_table" ADD COLUMN IF NOT EXISTS user_id uuid NOT NULL DEFAULT gen_random_uuid();`,
				`ALTER TABLE "public"."simple_table" DROP COLUMN IF EXISTS user_id;`,
			},
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing create and drlop column command #%d", n), func(t *testing.T) {
			drv := PgDriver{}
			sqlUp := (&drv).CreateColumnStatement(&tt.column)
			sqlDown := (&drv).DropColumnStatement(&tt.column)
			assert.Equal(t, sqlUp, tt.expected[0])
			assert.Equal(t, sqlDown, tt.expected[1])
		})
	}
}

func TestCreateDropIndexStatement(t *testing.T) {
	var testData = []struct {
		idx       IndexMeta
		excpected []string
	}{
		{
			IndexMeta{
				TableName: "auth_users",
				Name:      "idx_auth_users_id",
				Unique:    false,
				Type:      "btree",
				Where:     "",
				Option:    "",
				Parsed:    true,
				Columns: []IndexOption{
					IndexOption{
						Field:      "id",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   1,
					},
				},
			},
			[]string{
				`CREATE INDEX IF NOT EXISTS "idx_auth_users_id" ON "public"."auth_users" USING btree (id);`,
				`DROP INDEX IF EXISTS "idx_auth_users_id";`,
			},
		},
		{
			IndexMeta{
				TableName: "auth_users",
				Name:      "idx_auth_users_password",
				Unique:    true,
				Type:      "btree",
				Where:     "WHERE NOT user_type=5",
				Option:    "",
				Parsed:    true,
				Columns: []IndexOption{
					IndexOption{
						Field:      "password",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   1,
					},
				},
			},
			[]string{
				`CREATE UNIQUE INDEX IF NOT EXISTS "idx_auth_users_password" ON "public"."auth_users" USING btree (password);`,
				`DROP INDEX IF EXISTS "idx_auth_users_password";`,
			},
		},
		{
			IndexMeta{
				TableName: "auth_users",
				Name:      "idx_auth_users_password_upper",
				Unique:    true,
				Type:      "btree",
				Where:     "",
				Option:    "",
				Parsed:    true,
				Columns: []IndexOption{
					IndexOption{
						Field:      "password",
						Expression: "upper(password)",
						Sort:       "",
						Collate:    "",
						Priority:   1,
					},
				},
			},
			[]string{
				`CREATE UNIQUE INDEX IF NOT EXISTS "idx_auth_users_password_upper" ON "public"."auth_users" USING btree (upper(password));`,
				`DROP INDEX IF EXISTS "idx_auth_users_password_upper";`,
			},
		},
		{
			IndexMeta{
				TableName: "auth_users",
				Name:      "idx_auth_users_password_multi",
				Unique:    false,
				Type:      "btree",
				Where:     "",
				Option:    "",
				Parsed:    true,
				Columns: []IndexOption{
					IndexOption{
						Field:      "password",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   3,
					},
					IndexOption{
						Field:      "user_type",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   2,
					},
					IndexOption{
						Field:      "user_tag",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   4,
					},

					IndexOption{
						Field:      "user_id",
						Expression: "",
						Sort:       "",
						Collate:    "",
						Priority:   1,
					},
				},
			},
			[]string{
				`CREATE INDEX IF NOT EXISTS "idx_auth_users_password_multi" ON "public"."auth_users" USING btree (user_id, user_type, password, user_tag);`,
				`DROP INDEX IF EXISTS "idx_auth_users_password_multi";`,
			},
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing create index SQL statement #%d", n), func(t *testing.T) {
			drv := PgDriver{}
			sqlUp := (&drv).CreateIndexStatement(&tt.idx)
			sqlDown := (&drv).DropIndexStatement(&tt.idx)
			assert.Equal(t, sqlUp, tt.excpected[0])
			assert.Equal(t, sqlDown, tt.excpected[1])

		})
	}
}

func TestCreateDropConstraintStatement(t *testing.T) {
	var testData = []struct {
		ref       ReferenceMeta
		excpected []string
	}{
		{
			ReferenceMeta{
				TableName:  "commands",
				Name:       "commands_cars_car_id_id",
				Column:     "car_id",
				RefTable:   "cars",
				RefColumn:  "id",
				RefOptions: "ON DELETE CASCADE",
			},
			[]string{
				`ALTER TABLE "public"."commands" ADD CONSTRAINT "commands_cars_car_id_id" FOREIGN KEY (car_id) REFERENCES "public"."cars" (id) ON DELETE CASCADE;`,
				`ALTER TABLE "public"."commands" DROP CONSTRAINT IF EXISTS "commands_cars_car_id_id";`,
			},
		},
		{
			ReferenceMeta{
				TableName:  "auth_users",
				Name:       "fk_auth_users_users",
				Column:     "user_id",
				RefTable:   "users",
				RefColumn:  "id",
				RefOptions: "ON DELETE CASCADE ON UPDATE DO NOTHING",
			},
			[]string{
				`ALTER TABLE "public"."auth_users" ADD CONSTRAINT "fk_auth_users_users" FOREIGN KEY (user_id) REFERENCES "public"."users" (id) ON DELETE CASCADE ON UPDATE DO NOTHING;`,
				`ALTER TABLE "public"."auth_users" DROP CONSTRAINT IF EXISTS "fk_auth_users_users";`,
			},
		},
		{
			ReferenceMeta{
				TableName:  "commands",
				Name:       "fk_commands_cars_car_id_id_upd",
				Column:     "car_id",
				RefTable:   "cars",
				RefColumn:  "id",
				RefOptions: "ON DELETE RESTRICT",
			},
			[]string{
				`ALTER TABLE "public"."commands" ADD CONSTRAINT "fk_commands_cars_car_id_id_upd" FOREIGN KEY (car_id) REFERENCES "public"."cars" (id) ON DELETE RESTRICT;`,
				`ALTER TABLE "public"."commands" DROP CONSTRAINT IF EXISTS "fk_commands_cars_car_id_id_upd";`,
			},
		},
	}
	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing create constraint SQL statement #%d", n), func(t *testing.T) {
			drv := PgDriver{}
			sqlUp := (&drv).CreateConstraintStatement(&tt.ref)
			sqlDown := (&drv).DropConstraintStatement(&tt.ref)
			assert.Equal(t, sqlUp, tt.excpected[0])
			assert.Equal(t, sqlDown, tt.excpected[1])

		})
	}

}

package runner

import (
	"context"
	"fmt"
	"testing"

	utils "github.com/fastbear1/quack/internal"
	"github.com/stretchr/testify/assert"
)

func TestParseGormStructs(t *testing.T) {
	type Expected = struct {
		TableName     string
		Name          string
		Type          string
		IsNullable    bool
		ColumnDefault string
		IsPrimary     bool
	}
	type IdxColumn = struct {
		Field      string
		Expression string
		Priority   int
	}
	type ExpectedIndices = struct {
		TableName string
		Name      string
		Unique    bool
		Columns   []IdxColumn
		Type      string
		Where     string
	}

	var testData = []struct {
		data     ModelStruct
		expected []Expected
		expIdx   []ExpectedIndices
	}{
		{
			ModelStruct{
				Name: "TestTable",
				Fields: []FieldStruct{
					{
						FieldName: "ID",
						FieldType: "uuid.UUID",
						FieldTag:  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`,
					},
				},
			},
			[]Expected{
				{
					TableName:     "test_table",
					Name:          "id",
					Type:          "uuid",
					IsNullable:    false,
					ColumnDefault: "gen_random_uuid()",
					IsPrimary:     true,
				},
			},
			[]ExpectedIndices{},
		},
		{
			ModelStruct{
				Name: "TestTableA",
				Fields: []FieldStruct{
					{
						FieldName: "Name",
						FieldType: "string",
						FieldTag:  `gorm:"type:varchar(255);not null"`,
					},
					{
						FieldName: "SID",
						FieldType: "uint16",
						FieldTag:  `gorm:"index;not null"`,
					},
				},
			},
			[]Expected{
				{
					TableName:     "test_table_a",
					Name:          "name",
					Type:          "varchar(255)",
					IsNullable:    false,
					ColumnDefault: "",
					IsPrimary:     false,
				},
				{
					TableName:     "test_table_a",
					Name:          "sid",
					Type:          "smallint",
					IsNullable:    false,
					ColumnDefault: "",
					IsPrimary:     false,
				},
			},
			[]ExpectedIndices{
				{
					TableName: "test_table_a",
					Name:      "idx_test_table_a_sid",
					Unique:    false,
					Columns: []IdxColumn{
						{
							Field:      "sid",
							Expression: "",
							Priority:   10,
						},
					},
					Type:  "btree",
					Where: "",
				},
			},
		},
		{
			ModelStruct{
				Name: "TestTableb",
				Fields: []FieldStruct{
					{
						FieldName: "Name",
						FieldType: "string",
						FieldTag:  `gorm:"type:varchar(100);not null;default:active"`,
					},
					{
						FieldName: "SurName",
						FieldType: "string",
						FieldTag:  `gorm:"index:,expression:upper,unique;not null;default:undefined"`,
					},
					{
						FieldName: "SIDA",
						FieldType: "uint16",
						FieldTag:  `gorm:"index:idx_tableb_sidder,priority:5;not null"`,
					},
					{
						FieldName: "SIDB",
						FieldType: "int32",
						FieldTag:  `gorm:"index:idx_tableb_sidder,priority:7;not null"`,
					},
					{
						FieldName: "delta",
						FieldType: "float64",
						FieldTag:  `gorm:"index:idx_tableb_delta,,where:delta<0.5;not null"`,
					},
					{
						FieldName: "CreatedAt",
						FieldType: "time.Time",
						FieldTag:  `gorm:"type:timestamp;not null;default:now();\u003c-:create"`,
					},
				},
			},
			[]Expected{
				{
					TableName:     "test_tableb",
					Name:          "name",
					Type:          "varchar(100)",
					IsNullable:    false,
					ColumnDefault: `'active'::text`,
					IsPrimary:     false,
				},
				{
					TableName:     "test_tableb",
					Name:          "sur_name",
					Type:          "text",
					IsNullable:    false,
					ColumnDefault: `'undefined'::text`,
					IsPrimary:     false,
				},
				{
					TableName:     "test_tableb",
					Name:          "sida",
					Type:          "smallint",
					IsNullable:    false,
					ColumnDefault: "",
					IsPrimary:     false,
				},
				{
					TableName:     "test_tableb",
					Name:          "sidb",
					Type:          "bigint",
					IsNullable:    false,
					ColumnDefault: "",
					IsPrimary:     false,
				},
				{
					TableName:     "test_tableb",
					Name:          "delta",
					Type:          "double precision",
					IsNullable:    false,
					ColumnDefault: "",
					IsPrimary:     false,
				},
				{
					TableName:     "test_tableb",
					Name:          "created_at",
					Type:          "timestamp",
					IsNullable:    false,
					ColumnDefault: "now()",
					IsPrimary:     false,
				},
			},
			[]ExpectedIndices{
				{
					TableName: "test_tableb",
					Name:      "idx_test_tableb_sur_name_upper",
					Unique:    true,
					Columns: []IdxColumn{
						{
							Field:      "sur_name",
							Expression: "upper",
							Priority:   10,
						},
					},
					Type:  "btree",
					Where: "",
				},
				{
					TableName: "test_tableb",
					Name:      "idx_tableb_sidder",
					Unique:    false,
					Columns: []IdxColumn{
						{
							Field:      "sida",
							Expression: "",
							Priority:   5,
						},
						{
							Field:      "sidb",
							Expression: "",
							Priority:   7,
						},
					},
					Type:  "btree",
					Where: "",
				},
				{
					TableName: "test_tableb",
					Name:      "idx_tableb_delta",
					Unique:    false,
					Columns: []IdxColumn{
						{
							Field:      "delta",
							Expression: "",
							Priority:   10,
						},
					},
					Type:  "btree",
					Where: "delta<0.5",
				},
			},
		},
	}

	var conf utils.ConfigYaml
	conf.ReadConfig()
	conf.Database.Type = "postgres"

	ctx := context.Background()
	drv, _ := GetDriver(ctx, &conf)

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing parse model structs #%d\n", n), func(t *testing.T) {
			res := parseModelStruct(drv, tt.data)
			ColumnsAssertCount := 0
			// expected column data - cd
			for _, cd := range tt.expected {
				// resulted column data - rd
				for _, rd := range res.Columns {
					if cd.Name == rd.ColumnName {
						assert.Equal(t, rd.TableName, cd.TableName)
						assert.Equal(t, rd.DataType, cd.Type)
						assert.Equal(t, rd.IsNullable, cd.IsNullable)
						assert.Equal(t, rd.IsPrimary, cd.IsPrimary)
						assert.Equal(t, rd.ColumnDefault, cd.ColumnDefault)
						ColumnsAssertCount++
					}
				}
			}
			assert.Equal(t, ColumnsAssertCount, len(tt.expected))

			if len(tt.expIdx) > 0 {
				IndicesAssertCount := 0
				// expected indices - id
				for _, id := range tt.expIdx {
					IndicesAsetColumnCount := 0
					// resulted indices data - ri
					for _, ri := range res.Indeces {
						if ri.Name == id.Name {
							assert.Equal(t, ri.Type, id.Type)
							assert.Equal(t, ri.Unique, id.Unique)
							assert.Equal(t, ri.Where, id.Where)
							IndicesAssertCount++
							// expected indices columns - idc\
							for _, idc := range id.Columns {
								// resulted indices columns
								for _, ric := range ri.Columns {
									if ric.Field == idc.Field {
										assert.Equal(t, ric.Priority, idc.Priority)
										assert.Equal(t, ric.Expression, idc.Expression)
										IndicesAsetColumnCount++
									}
								}
							}
							assert.Equal(t, IndicesAsetColumnCount, len(id.Columns))
						}
					}
				}
				assert.Equal(t, IndicesAssertCount, len(tt.expIdx))
			}
		})
	}
}

func TestParseReferenceEmbedStructs(t *testing.T) {
	type Result struct {
		Name      string
		Column    string
		RefTable  string
		RefColumn string
		Options   string
	}

	var testDate = []struct {
		tableName string
		refTable  string
		tag       string
		expected  Result
	}{
		{
			tableName: "TestTable",
			refTable:  "Users",
			tag:       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:SET NULL;"`,
			expected: Result{
				Name:      "fk_TestTable_user_id__users_id",
				Column:    "user_id",
				RefTable:  "users",
				RefColumn: "id",
				Options:   "ON DELETE CASCADE ON UPDATE SET NULL",
			},
		},
		{
			tableName: "TestTable",
			refTable:  "Cars",
			tag:       `gorm:"foreignKey:CarId;referenceName:commands_cars_car_id_id;constraint:OnDelete:CASCADE;"`,
			expected: Result{
				Name:      "commands_cars_car_id_id",
				Column:    "car_id",
				RefTable:  "cars",
				RefColumn: "id",
				Options:   "ON DELETE CASCADE",
			},
		},
		{
			tableName: "TestTable",
			refTable:  "Owners",
			tag:       `gorm:"foreignKey:OwnerId;referenceName:commands_owner_owner_id_id;constraint:OnDelete:CASCADE;"`,
			expected: Result{
				Name:      "commands_owner_owner_id_id",
				Column:    "owner_id",
				RefTable:  "owners",
				RefColumn: "id",
				Options:   "ON DELETE CASCADE",
			},
		},
	}

	var conf utils.ConfigYaml
	conf.ReadConfig()
	conf.Database.Type = "postgres"

	ctx := context.Background()
	drv, _ := GetDriver(ctx, &conf)

	for n, tt := range testDate {
		t.Run(fmt.Sprintf("Test for parsing reference column #%d", n), func(t *testing.T) {
			res := parseReferenceEmbedStructs(drv, tt.tableName, tt.refTable, tt.tag)
			assert.Equal(t, res.Name, tt.expected.Name)
			assert.Equal(t, res.Column, tt.expected.Column)
			assert.Equal(t, res.RefTable, tt.expected.RefTable)
			assert.Equal(t, res.RefColumn, tt.expected.RefColumn)
			assert.Equal(t, res.RefOptions, tt.expected.Options)

		})
	}
}

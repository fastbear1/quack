package runner

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	. "github.com/fastbear1/quack/schema"
	"github.com/stretchr/testify/assert"
)

func TestIsColumnSchemaChanged(t *testing.T) {
	var testData = []struct {
		ColumnA  Column
		ColumnB  Column
		expected bool
	}{
		{
			Column{
				ColumnName:    "ColumnA",
				DataType:      "string",
				IsNullable:    false,
				ColumnDefault: "",
			},
			Column{
				ColumnName:    "ColumnA",
				DataType:      "string",
				IsNullable:    false,
				ColumnDefault: "",
			},
			false,
		},
		{
			Column{
				ColumnName:    "ColumnA",
				DataType:      "string",
				IsNullable:    false,
				ColumnDefault: "",
			},
			Column{
				ColumnName:    "ColumnA",
				DataType:      "string",
				IsNullable:    false,
				ColumnDefault: `"active"::text`,
			},
			true,
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing column schema changed #%d\n", n), func(t *testing.T) {
			res := isColumnSchemaChanged(&tt.ColumnA, &tt.ColumnB)
			assert.Equal(t, res, tt.expected)
		})
	}
}

func TestIsReferenceSchemaChanged(t *testing.T) {
	var testData = []struct {
		ColumnA  ReferenceMeta
		ColumnB  ReferenceMeta
		expected bool
	}{
		{
			ReferenceMeta{
				Name:       "ref_table_user_ud_users",
				RefColumn:  "user_id",
				RefTable:   "users",
				RefOptions: "ON DELETE CASCADE",
			},
			ReferenceMeta{
				Name:       "ref_table_user_ud_users",
				RefColumn:  "user_id",
				RefTable:   "users",
				RefOptions: "ON DELETE CASCADE",
			},
			false,
		},
		{
			ReferenceMeta{
				Name:       "ref_table_user_ud_users",
				RefColumn:  "user_id",
				RefTable:   "users",
				RefOptions: "ON DELETE CASCADE",
			},
			ReferenceMeta{
				Name:       "ref_table_user_ud_users",
				RefColumn:  "user_id",
				RefTable:   "users",
				RefOptions: "",
			},
			true,
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing column schema changed #%d\n", n), func(t *testing.T) {
			res := isReferenceSchemaChanged(&tt.ColumnA, &tt.ColumnB)
			assert.Equal(t, res, tt.expected)
		})
	}
}

func TestIsIndexSchemaChanged(t *testing.T) {
	var testData = []struct {
		ColumnA  IndexMeta
		ColumnB  IndexMeta
		expected bool
	}{
		{
			IndexMeta{
				Name:   "idx_table_column",
				Unique: false,
				Type:   "btree",
				Columns: []IndexOption{
					{
						Field:      "num",
						Expression: "sum",
						Priority:   1,
					},
					{
						Field:      "delta",
						Expression: "sum",
						Priority:   2,
					},
				},
			},
			IndexMeta{
				Name:   "idx_table_column",
				Unique: false,
				Type:   "btree",
				Columns: []IndexOption{
					{
						Field:      "num",
						Expression: "sum",
						Priority:   1,
					},
					{
						Field:      "delta",
						Expression: "sum",
						Priority:   2,
					},
				},
			},
			false,
		},
		{
			IndexMeta{
				Name:   "idx_table_column",
				Unique: false,
				Type:   "btree",
				Columns: []IndexOption{
					{
						Field:      "num",
						Expression: "sum",
						Priority:   2,
					},
					{
						Field:      "delta",
						Expression: "sum",
						Priority:   1,
					},
				},
			},
			IndexMeta{
				Name:   "idx_table_column",
				Unique: false,
				Type:   "btree",
				Columns: []IndexOption{
					{
						Field:      "num",
						Expression: "sum",
						Priority:   1,
					},
					{
						Field:      "delta",
						Expression: "sum",
						Priority:   2,
					},
				},
			},
			true,
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing column schema changed #%d\n", n), func(t *testing.T) {
			res := isIndexSchemaChanged(&tt.ColumnA, &tt.ColumnB)
			assert.Equal(t, res, tt.expected)
		})
	}
}

func TestCompareMetaState(t *testing.T) {
	var testData = []struct {
		leftData  []TableMeta
		rightData []TableMeta
		expUp     []string
		expDown   []string
	}{
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(255)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]TableMeta{},
			[]string{"CreateTable"},
			[]string{"DeleteTable"},
		},
		{
			[]TableMeta{},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(255)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]string{"DeleteTable"},
			[]string{"CreateTable"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(255)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
					},
				},
			},
			[]string{"CreateColumn"},
			[]string{"DeleteColumn"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(255)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]string{"DeleteColumn"},
			[]string{"CreateColumn"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(255)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "name",
							TableName:     "TestTable",
							DataType:      "varchar(100)",
							IsNullable:    true,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]string{"AlterColumn"},
			[]string{"AlterColumn"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					References: []ReferenceMeta{
						{
							TableName:  "TestTable",
							Name:       "ref_test_table_user_id",
							Column:     "user_id",
							RefTable:   "users",
							RefColumn:  "id",
							RefOptions: "ON DELETE CASCADE",
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]string{"CreateConstraint"},
			[]string{"DeleteConstraint"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					References: []ReferenceMeta{
						{
							TableName:  "TestTable",
							Name:       "ref_test_table_user_id",
							Column:     "user_id",
							RefTable:   "users",
							RefColumn:  "id",
							RefOptions: "ON DELETE CASCADE",
						},
					},
				},
			},
			[]string{"DeleteConstraint"},
			[]string{"CreateConstraint"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					References: []ReferenceMeta{
						{
							TableName:  "TestTable",
							Name:       "ref_test_table_user_id",
							Column:     "user_id",
							RefTable:   "users",
							RefColumn:  "id",
							RefOptions: "ON DELETE CASCADE",
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "user_id",
							TableName:     "TestTable",
							DataType:      "bigint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					References: []ReferenceMeta{
						{
							TableName:  "TestTable",
							Name:       "ref_test_table_user_id",
							Column:     "user_id",
							RefTable:   "users",
							RefColumn:  "id",
							RefOptions: "",
						},
					},
				},
			},
			[]string{"DeleteConstraint", "CreateConstraint"},
			[]string{"DeleteConstraint", "CreateConstraint"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					Indeces: []IndexMeta{
						{
							TableName: "TestTable",
							Name:      "idx_test_table_sid_index",
							Unique:    true,
							Type:      "btree",
							Columns: []IndexOption{
								{
									Field:      "sid",
									Expression: "",
									Priority:   1,
								},
							},
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]string{"CreateIndex"},
			[]string{"DropIndex"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					Indeces: []IndexMeta{
						{
							TableName: "TestTable",
							Name:      "idx_test_table_sid_index",
							Unique:    true,
							Type:      "btree",
							Columns: []IndexOption{
								{
									Field:      "sid",
									Expression: "",
									Priority:   1,
								},
							},
						},
					},
				},
			},
			[]string{"DropIndex"},
			[]string{"CreateIndex"},
		},
		{
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					Indeces: []IndexMeta{
						{
							TableName: "TestTable",
							Name:      "idx_test_table_sid_index",
							Unique:    true,
							Type:      "btree",
							Columns: []IndexOption{
								{
									Field:      "sid",
									Expression: "",
									Priority:   1,
								},
							},
						},
					},
				},
			},
			[]TableMeta{
				{
					Name: "TestTable",
					Columns: []Column{
						{
							ColumnName:    "id",
							TableName:     "TestTable",
							DataType:      "uuid",
							IsNullable:    false,
							ColumnDefault: "uuidv4()",
							IsPrimary:     true,
						},
						{
							ColumnName:    "sid",
							TableName:     "TestTable",
							DataType:      "smallint",
							IsNullable:    false,
							ColumnDefault: "",
							IsPrimary:     false,
						},
					},
					Indeces: []IndexMeta{
						{
							TableName: "TestTable",
							Name:      "idx_test_table_sid_index",
							Unique:    false,
							Type:      "btree",
							Columns: []IndexOption{
								{
									Field:      "sid",
									Expression: "",
									Priority:   1,
								},
							},
						},
					},
				},
			},
			[]string{"DropIndex", "CreateIndex"},
			[]string{"DropIndex", "CreateIndex"},
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing compare meta state (%s) #%d\n", tt.expUp[0], n), func(t *testing.T) {
			fUp, fDown := compareMetaState(tt.leftData, tt.rightData)
			for n, fu := range fUp {
				funcId := runtime.FuncForPC(reflect.ValueOf(fu).Pointer()).Name()
				funcName := strings.TrimSuffix(strings.Split(funcId, ".")[3], "-fm")
				assert.Equal(t, funcName, tt.expUp[n])
			}
			for n, fd := range fDown {
				funcId := runtime.FuncForPC(reflect.ValueOf(fd).Pointer()).Name()
				funcName := strings.TrimSuffix(strings.Split(funcId, ".")[3], "-fm")
				assert.Equal(t, funcName, tt.expDown[n])
			}
		})
	}
}

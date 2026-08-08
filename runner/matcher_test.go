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
		expected  string
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
			"CreateTable",
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
			"DeleteTable",
		},
	}

	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing compare meta state #%d\n", n), func(t *testing.T) {
			fUp, _ := compareMetaState(tt.leftData, tt.rightData)
			funcId := runtime.FuncForPC(reflect.ValueOf(fUp[0]).Pointer()).Name()
			funcName := strings.TrimSuffix(strings.Split(funcId, ".")[3], "-fm")
			fmt.Println(runtime.FuncForPC(reflect.ValueOf(fUp[0]).Pointer()).Name())

			assert.Equal(t, funcName, tt.expected)
		})
	}
}

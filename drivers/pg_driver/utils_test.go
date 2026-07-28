package pgdriver

import (
	"fmt"
	"testing"

	. "github.com/fastbear1/quack/schema"
	"github.com/stretchr/testify/assert"
)

func TestParseDatabaseIndices(t *testing.T) {
	type ExcData struct {
		err error
		idx IndexMeta
	}

	var testData = []struct {
		indexDef string
		excepted ExcData
	}{
		{"CREATE INDEX idx_auth_users_id ON public.auth_users USING btree (id)",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
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
			},
		},
		{"CREATE INDEX idx_auth_users_password ON public.auth_users USING btree (password)",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password",
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
							Priority:   1,
						},
					},
				},
			},
		},
		{"CREATE UNIQUE INDEX idx_auth_users_password_user_id_2 ON public.auth_users USING btree (password, user_id)",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password_user_id_2",
					Unique:    true,
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
							Priority:   1,
						},
						IndexOption{
							Field:      "user_id",
							Expression: "",
							Sort:       "",
							Collate:    "",
							Priority:   2,
						},
					},
				},
			},
		},
		{"CREATE INDEX idx_auth_users_password_lower ON public.auth_users USING btree (upper(password))",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password_lower",
					Unique:    false,
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
			},
		},
		{"CREATE INDEX idx_auth_users_password_partial_id ON public.auth_users USING btree (password) WHERE (NOT ((created_at)::date > '2026-01-01'::date))",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password_partial_id",
					Unique:    false,
					Type:      "btree",
					Where:     " WHERE (NOT ((created_at)::date > '2026-01-01'::date))",
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
			},
		},
		{"CREATE UNIQUE INDEX idx_auth_users_password_user_id ON public.auth_users USING btree (password) INCLUDE (user_id)",
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password_user_id",
					Unique:    true,
					Type:      "btree",
					Where:     "",
					Option:    " INCLUDE (user_id)",
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
			},
		},
		{`CREATE INDEX idx_auth_users_password_collation ON public.auth_users USING btree (password COLLATE "de-DE-x-icu")`,
			ExcData{
				err: nil,
				idx: IndexMeta{
					TableName: "public.auth_users",
					Name:      "idx_auth_users_password_collation",
					Unique:    false,
					Type:      "btree",
					Where:     "",
					Option:    "",
					Parsed:    true,
					Columns: []IndexOption{
						IndexOption{
							Field:      "password ",
							Expression: "",
							Sort:       "",
							Collate:    `COLLATE "de-DE-x-icu"`,
							Priority:   1,
						},
					},
				},
			},
		},
	}
	for n, tt := range testData {
		t.Run(fmt.Sprintf("Testing parse index definition #%d", n), func(t *testing.T) {
			res, err := ParseDatabaseIndices(tt.indexDef)
			assert.Nil(t, err)
			assert.Equal(t, res.Parsed, tt.excepted.idx.Parsed)
			assert.Equal(t, res.Name, tt.excepted.idx.Name)
			assert.Equal(t, res.Unique, tt.excepted.idx.Unique)
			assert.Equal(t, res.Type, tt.excepted.idx.Type)
			assert.Equal(t, res.Where, tt.excepted.idx.Where)
			assert.Equal(t, res.Option, tt.excepted.idx.Option)
			assert.Equal(t, res.Parsed, tt.excepted.idx.Parsed)
			assert.Equal(t, len(res.Columns), len(tt.excepted.idx.Columns))
			for i, f := range res.Columns {
				assert.Equal(t, f.Field, tt.excepted.idx.Columns[i].Field)
				assert.Equal(t, f.Expression, tt.excepted.idx.Columns[i].Expression)
				assert.Equal(t, f.Sort, tt.excepted.idx.Columns[i].Sort)
				assert.Equal(t, f.Collate, tt.excepted.idx.Columns[i].Collate)
				assert.Equal(t, f.Priority, tt.excepted.idx.Columns[i].Priority)
			}
		})
	}
}

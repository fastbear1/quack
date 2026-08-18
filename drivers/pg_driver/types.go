package pgdriver

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Pg driver type definitions
type PgColumn struct {
	Column_name              string
	Data_type                string
	Character_maximum_length pgtype.Uint32
	Is_nullable              string
	Column_default           pgtype.Text
	Is_identity              string
	Identity_generation      pgtype.Text
}

type AlterData struct {
	TableName     string
	ColumnName    string
	Type          uint8
	DataType      string
	IsNullable    bool
	ColumnDefault string
}

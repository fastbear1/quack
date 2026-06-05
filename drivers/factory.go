package drivers

import (
	"errors"

	"github.com/fastbear1/quack/internal/utils"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type DbHandler interface {
	GetData(conf *utils.ConfigYaml) string
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
		case "postgres":
			return PgHandler, nil
 		default:
			return nil, ErrNotFound
	}
}

type PgHandler struct {}

func (pg *PgHandler) GetData(conf *utils.ConfigYaml) string {
	return string
}


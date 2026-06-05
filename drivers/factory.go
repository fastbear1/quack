package drivers

import (
	"errors"

	utils "github.com/fastbear1/quack/internal"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

type DbHandler interface {
	GetData(conf *utils.ConfigYaml) ([]string, error)
}

func GetDriver(db_type string) (DbHandler, error) {
	switch db_type {
	case "postgres":
		return &PgHandler{}, nil
	default:
		return nil, ErrNotFound
	}
}

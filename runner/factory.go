package runner

import (
	"context"
	"errors"

	pg "github.com/fastbear1/quack/drivers/pg_driver"
	utils "github.com/fastbear1/quack/internal"
	. "github.com/fastbear1/quack/schema"
)

var ErrNotFound = errors.New("resource not found. Can't find proper database handler")

func GetDriver(ctx context.Context, conf *utils.ConfigYaml) (DbInterface, error) {
	switch conf.Database.Type {
	case "postgres":
		return pg.GetPgDriver(ctx, conf), nil
	default:
		return nil, ErrNotFound
	}
}

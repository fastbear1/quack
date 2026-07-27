package runner

import (
	"context"
	"testing"

	pg "github.com/fastbear1/quack/drivers/pg_driver"
	utils "github.com/fastbear1/quack/internal"
	. "github.com/fastbear1/quack/schema"
)

func TestFactoryMethod(t *testing.T) {
	var drv DbInterface
	var ctx context.Context = context.Background()
	var conf utils.ConfigYaml
	conf.ReadConfig()
	conf.Database.Type = "postgres"
	drv, err := GetDriver(ctx, &conf)
	if err != nil {
		t.Error()
	}
	if _, ok := drv.(*pg.PgDriver); ok != true {
		t.Error()
	}
}

func TestFactoryMethodUnknowHandler(t *testing.T) {
	var ctx context.Context = context.Background()
	var conf utils.ConfigYaml
	conf.ReadConfig()
	conf.Database.Type = "not-postgres"
	_, err := GetDriver(ctx, &conf)
	if err == nil {
		t.Error()
	}
}

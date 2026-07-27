package runner

import (
	"testing"
	sch "github.com/fastbear1/quack/schema"
	pg "github.com/fastbear1/quack/drivers/pg_driver"
)

func TestFactoryMethod(t *testing.T) {
	var drv sch.DbHandler
	drv, err := GetDriver("postgres")
	if err != nil {
		t.Error()
	}
	if _, ok := drv.(*pg.PgHandler); ok != true {
		t.Error()
	}
}

func TestFactoryMethodUnknowHandler(t *testing.T) {
	_, err := GetDriver("not-postgres")
	if err == nil {
		t.Error()
	}
}

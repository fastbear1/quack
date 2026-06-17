package runner

import (
	"testing"

	pg "github.com/fastbear1/quack/drivers/postgres"
)

func TestFactoryMethod(t *testing.T) {
	var drv DbHandler
	drv, err := GetDriver("postgres")
	if err != nil {
		t.Error()
	}
	if _, ok := drv.(pg.PgHandler); ok != nil {
		t.Error()
	}
}

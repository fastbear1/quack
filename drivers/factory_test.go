package drivers

import (
	"testing"
)

func TestFactoryMethod(t *testing.T) {
	var drv DbHandler
	drv, err := GetDriver("postgres")
	if err != nil {
		t.Error()
	}
	if _, ok := drv.(*PgHandler); ok != true {
		t.Error()
	}
}

func TestFactoryMethodUnknowHandler(t *testing.T) {
	_, err := GetDriver("not-postgres")
	if err == nil {
		t.Error()
	}
}

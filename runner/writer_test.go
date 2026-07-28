package runner

import (
	"os"
	"strings"
	"testing"

	utils "github.com/fastbear1/quack/internal"
	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	var conf utils.ConfigYaml
	conf.ReadConfig()
	conf.Migrations.Path = "migrations"

	err := os.Mkdir(conf.Migrations.Path.String(), 0777)
	assert.Nil(t, err)
	//defer utils.CleanUpDir(conf.Migrations.Path.String())

	writeToFile(&conf, "test", []string{"one", "two", "three"}, []string{"one", "two", "three"})
	files, err := os.ReadDir(conf.Migrations.Path.String())
	assert.Nil(t, err)

	for _, fl := range files {
		assert.True(t, strings.Contains(fl.Name(), "_test.sql"))
	}
}

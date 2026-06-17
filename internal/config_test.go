package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConfigFile(t *testing.T) {
	_, err := FindConfigFile()
	assert.Nil(t, err)
}

func TestReadConfigFile(t *testing.T) {
	conf := ConfigYaml{}
	conf.ReadConfig()

	assert.Equal(t, conf.Version, float32(0.1))
	assert.Equal(t, conf.Database.Uri.String(), "postgres")
	assert.Equal(t, conf.Database.Name.String(), "test")
	assert.Equal(t, conf.Models.Path.String(), "test")
	assert.Equal(t, conf.Migrations.Path.String(), "test")

	// assert lists
	var modelExcld, dbexcld []string
	for _, v := range conf.Models.Exclude {
		modelExcld = append(modelExcld, string(v))
	}
	for _, v := range conf.Database.Exclude {
		dbexcld = append(dbexcld, string(v))
	}
	assert.Equal(t, modelExcld, []string{"Base", "Test"})
	assert.Equal(t, dbexcld, []string{"test_table_1"})

}

func TestStringVal(t *testing.T) {
	var str StringVal
	err := str.Set("empty")
	assert.Nil(t, err)
	assert.Equal(t, str.String(), "empty")
}

func TestStringList(t *testing.T) {
	var strl StringList
	err := strl.Set("key,val")
	assert.Nil(t, err)
	assert.Equal(t, strl.String(), "key,val")
}

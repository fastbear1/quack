package utils

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
	var sarr = []string{"one", "two", "free"}
	assert.True(t, InArray(sarr, "one"))
	assert.False(t, InArray(sarr, "five"))
}

func TestCleanUpDirectory(t *testing.T) {
	dirName := "test_test"
	fileName := "file.txt"
	err := os.Mkdir(dirName, 0777)
	assert.Nil(t, err)

	file, err := os.Create(path.Join(dirName, fileName))
	defer file.Close()
	assert.Nil(t, err)

	CleanUpDir(dirName)

	_, ok := os.Stat(dirName)
	assert.NotNil(t, ok)

}

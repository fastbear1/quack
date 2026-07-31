package utils

import (
	"fmt"
	"math/rand"
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

func TestSortArray(t *testing.T) {
	type Cell struct {
		Name     string
		Position int
	}
	var data = make([][]Cell, 10)

	ln := 10
	for i := range 10 {
		intCell := make([]Cell, ln)
		for j := range ln {
			intCell[j] = Cell{
				Name:     getRandomName(),
				Position: getRandomInt(),
			}
		}
		data[i] = intCell
		ln = ln + 10
	}

	for n, tt := range data {
		t.Run(fmt.Sprintf("Testing insertion sort algorithm #%d", n), func(t *testing.T) {
			SortArray(tt, func(i, j int) bool {
				return tt[i].Position > tt[j].Position
			})
			for n, vl := range tt {
				if n < len(tt)-2 && vl.Position > tt[n+1].Position {
					t.Error()
				}
			}
		})
	}
}

func getRandomName() string {
	a := 65 + rand.Intn(25)
	b := 65 + rand.Intn(25)
	c := 65 + rand.Intn(25)
	return fmt.Sprintf("%c%c%c", a, b, c)
}

func getRandomInt() int {
	return rand.Intn(100)
}

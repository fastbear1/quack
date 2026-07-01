package runner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCatalogData(t *testing.T) {

	var test = []struct {
		name        string
		left        []string
		right       []string
		expectLeft  []string
		expectRight []string
	}{
		{
			"Testing compare array elements #1",
			[]string{"one", "two", "three", "five"},
			[]string{"one", "two", "four", "five"},
			[]string{"four"},
			[]string{"three"},
		},
		{
			"Testing compare array elements #2",
			[]string{"ALTER", "TABLE", "SET", "DEFAULT", "not", "Null"},
			[]string{"ALTER", "TABLE", "ADD", "Null", "DEFAULT", "SET"},
			[]string{"ADD"},
			[]string{"not"},
		},
	}
	for _, tt := range test {
		t.Run(fmt.Sprintf("Test for %s", tt.name), func(t *testing.T) {
			resLeft, resRight := getCatalogData(tt.left, tt.right)
			assert.Equal(t, resLeft, tt.expectLeft)
			assert.Equal(t, resRight, tt.expectRight)
		})
	}
}

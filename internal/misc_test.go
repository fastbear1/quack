package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
	var sarr = []string{"one", "two", "free"}
	assert.True(t, InArray(sarr, "one"))
	assert.False(t, InArray(sarr, "five"))
}

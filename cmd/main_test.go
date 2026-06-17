package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFalgs(t *testing.T) {
	var tt = struct {
		flagStr  string
		excepted []string
	}{
		"--models=models --path=migrations --uri=postgres://user:pass@host:port/database",
		[]string{"models", "migrations", "postgres://user:pass@host:port/database", "Base,Test"},
	}
	flagsCmd := strings.Split(tt.flagStr, " ")
	for _, f := range flagsCmd {
		os.Args = append(os.Args, f)
	}
	conf := ParseFlags()
	assert.Equal(t, conf.Models.Path.String(), tt.excepted[0])
	assert.Equal(t, conf.Migrations.Path.String(), tt.excepted[1])
	assert.Equal(t, conf.Database.Uri.String(), tt.excepted[2])

	var exclude []string
	for _, v := range conf.Models.Exclude {
		exclude = append(exclude, v)
	}
	assert.Equal(t, exclude, strings.Split(tt.excepted[3], ","))
}

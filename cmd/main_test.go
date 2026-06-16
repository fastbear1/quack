package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFalgs(t *testing.T) {
	var tests = []struct {
		name     string
		flagStr  string
		excepted []string
	}{
		{
			"testing flag args parser",
			"--models=models --path=migrations --uri=postgres://user:pass@host:port/database --exclude=Base,TestUsers",
			[]string{"models", "migrations", "postgres://user:pass@host:port/database", "Base,TestUsers"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagsCmd := strings.Split(tt.flagStr, " ")
			for _, f := range flagsCmd {
				os.Args = append(os.Args, f)
			}
			conf := ParseFlags()
			assert.Equal(t, conf.Models.Path.String(), tt.excepted[0])
		})
	}
}

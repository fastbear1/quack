package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	utils "github.com/fastbear1/quack/internal"
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

func TestValidateConfig(t *testing.T) {
	var test = []struct {
		Uri      string
		Type     string
		excepted bool
	}{
		{"postgres://user:pass@host:port/database", "postgres", true},
		{"postgres//user_pass@host_port/database", "", false},
		{"postg:res://user:pass@host:port/database", "postg", true},
	}
	var conf utils.ConfigYaml
	conf.ReadConfig()

	for n, tt := range test {
		t.Run(fmt.Sprintf("Test for validate config #%d", n), func(t *testing.T) {
			conf.Database.Uri = utils.StringVal(tt.Uri)
			res := isConfigValid(&conf)
			fmt.Println(res)
			assert.Equal(t, res, tt.excepted)
			if res {
				assert.Equal(t, conf.Database.Type, tt.Type)
			}
		})
	}

}

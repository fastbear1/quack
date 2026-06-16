package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseFalgs(t *testing.T) {
	flagsStr := "--models=models --path=migrations --uri=postgres://user:pass@host:port/database --exclude=Base,TestUsers"

	flagsCmd := strings.Split(flagsStr, " ")

	fmt.Println(os.Args)
	for _, f := range flagsCmd {
		os.Args = append(os.Args, f)
	}

	conf := ParseFlags()
	fmt.Println(conf)
}

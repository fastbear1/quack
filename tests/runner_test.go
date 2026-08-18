package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	utils "github.com/fastbear1/quack/internal"
	proc "github.com/fastbear1/quack/runner"
	"github.com/stretchr/testify/assert"
)

func TestInitDatabase(t *testing.T) {
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case1/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-init-database")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

	//getting migration file name
	cmdres, _ := exec.Command("ls", "./migrations/").Output()
	// check that content of migration file is identical to original
	filename := strings.TrimSuffix(string(cmdres), "\n")
	assert.True(t, CompareFiles("../sandbox/case1/migrations/init-database.sql", fmt.Sprintf("./migrations/%s", filename)))

	home := os.Getenv("HOME")
	cmd := exec.Command(fmt.Sprintf("%s/go/bin/goose", home), "postgres", "user=quack dbname=quack password=pass host=postgres sslmode=disable", "-dir=migrations", "up")

	var out bytes.Buffer
	cmd.Stderr = &out

	excErr := cmd.Run()
	assert.Nil(t, excErr)
	if excErr != nil {
		fmt.Println(out.String())
	}
}

func getTestConfig() *utils.ConfigYaml {
	var conf utils.ConfigYaml

	conf.ReadConfig()
	conf.Database.Uri = "postgres://quack:pass@postgres:5432/quack"
	conf.Database.Name = "quack"
	conf.Database.Type = "postgres"
	conf.Database.Exclude = append(conf.Database.Exclude, "goose_migrations", "goose_db_version")

	return &conf
}

func CompareFiles(file1, file2 string) bool {
	chunkSize := 4 * 1024

	// compare contents
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	b1 := make([]byte, chunkSize)
	b2 := make([]byte, chunkSize)
	for {
		n1, err1 := io.ReadFull(f1, b1)
		n2, err2 := io.ReadFull(f2, b2)

		if !bytes.Equal(b1[:n1], b2[:n2]) {
			return false
		}

		if (err1 == io.EOF && err2 == io.EOF) || (err1 == io.ErrUnexpectedEOF && err2 == io.ErrUnexpectedEOF) {
			return true
		}

		if err1 != nil || err2 != nil {
			return false
		}
	}
}

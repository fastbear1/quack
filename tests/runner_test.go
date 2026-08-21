package tests

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

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
	lsc, _ := exec.Command("ls", `./migrations/`).Output()
	lscfile := strings.Split(string(lsc), "\n")[0]
	// check that content of migration file is identical to original
	assert.True(t, CompareFiles("../sandbox/case1/migrations/init-database.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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

func TestManageTableColumns(t *testing.T) {
	//t.Skip("Skipping test")
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case2/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-manage-columns")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

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

func TestAlterTableColumns(t *testing.T) {
	//t.Skip("Skipping test")
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case3/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-alter-columns")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

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

func TestTableIndices(t *testing.T) {
	//t.Skip("Skipping test")
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case4/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-table-indices")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

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

func TestTableConstraints(t *testing.T) {
	//t.Skip("Skipping test")
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case5/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-table-constraints")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

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

package tests

import (
	"bufio"
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
	assert.True(t, LinearFilesComparison("../sandbox/case1/migrations/init-database.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case2/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-manage-columns")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

	//getting migration file name
	lsc, _ := exec.Command("ls", `./migrations/`).Output()
	lscfile := strings.Split(string(lsc), "\n")[1]
	// check that content of migration file is identical to original
	assert.True(t, LinearFilesComparison("../sandbox/case2/migrations/manage-columns.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case3/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-alter-columns")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

	//getting migration file name
	lsc, _ := exec.Command("ls", `./migrations/`).Output()
	lscfile := strings.Split(string(lsc), "\n")[2]
	// check that content of migration file is identical to original
	assert.True(t, LinearFilesComparison("../sandbox/case3/migrations/alter-columns.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case4/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-table-indices")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

	//getting migration file name
	lsc, _ := exec.Command("ls", `./migrations/`).Output()
	lscfile := strings.Split(string(lsc), "\n")[3]
	// check that content of migration file is identical to original
	assert.True(t, LinearFilesComparison("../sandbox/case4/migrations/table-indices.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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
	time.Sleep(time.Second)
	conf := getTestConfig()
	conf.Models.Path = "../sandbox/case5/models"
	conf.Migrations.Path = "migrations"

	ctx := context.Background()
	migrationFile := fmt.Sprintf("test-table-constraints")

	err := proc.Run(ctx, conf, migrationFile)
	assert.Equal(t, err, 0)

	//getting migration file name
	lsc, _ := exec.Command("ls", `./migrations/`).Output()
	lscfile := strings.Split(string(lsc), "\n")[4]
	// check that content of migration file is identical to original
	assert.True(t, LinearFilesComparison("../sandbox/case5/migrations/table-constraints.sql", fmt.Sprintf("./migrations/%s", lscfile)))

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
	conf.Verbose = false

	return &conf
}

func LinearFilesComparison(leftFile string, rightFile string) bool {
	f1, err := os.Open(leftFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f1.Close()

	f2, err := os.Open(rightFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f2.Close()

	bf1 := make([][]byte, 0)
	bf2 := make([][]byte, 0)

	scanner := bufio.NewScanner(f1)
	for scanner.Scan() {
		bf1 = append(bf1, scanner.Bytes())
	}

	scanner = bufio.NewScanner(f2)
	for scanner.Scan() {
		bf2 = append(bf2, scanner.Bytes())
	}

	var cutline int
	bff := bf2[:len(bf2)]
	for _, line := range bf1 {
		for i, ln := range bff {
			if bytes.Equal(line, ln) {
				cutline = i
				break
			}
		}
		bff = append(bff[0:cutline], bff[cutline+1:]...)
	}

	if len(bff) == 0 {
		return true
	}
	return false
}

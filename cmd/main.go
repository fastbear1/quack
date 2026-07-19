package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	utils "github.com/fastbear1/quack/internal"
	proc "github.com/fastbear1/quack/internal/runner"
)

const version string = "0.23.2"

const (
	helpInfo = `Quack - generate migration file for goose according gorm struct models 
information and database state. Use config file quack_config.yaml for running params.`
	usageInfo = `
Usage:
  quack [flags] command filename[optional]

filename is optional, if not provided default name will be used for migration file.
Commands:
  - run - quack(run and create) goose migration file
  - help - show help information
  - version - show current version(can be used for checking config file)
  command usage example:
    - quack help - show help information
    - quack run - run creating a new migration files
flags:
  - models - directory where all gorm struct models live
  - uri - URI connection database string
  - dbname - database name
  - path - directory where a new migration files will be stored
  - exclude - exclude gorm struct models(usually it's a embeded struct or not a model struct)
  - db-exclude - exclude database tables(for example goose_migrations table)
  flag usage examples:
    - quack --models=models --path=migrations --uri=postgres://user:pass@host:port/database --exclude=Base,TestUsers run
    - quack --models=internal/models --path=models --uri=postgres://user:pass@host:port/database --exclude=Base --db-exclude=goose,goose_migrations run`
)

func main() {
	conf := ParseFlags()
	commands := flag.Args()
	if len(commands) > 0 {
		switch commands[0] {
		case "help":
			fmt.Println(helpInfo)
			fmt.Println(usageInfo)
		case "run":
			fmt.Println("Quacking migration file")
			var fileName string
			if len(commands) > 1 {
				fileName = strings.ToLower(commands[1])
			} else {
				fmt.Println("Filename not provided. Using default name 'goose_file'")
				fileName = "goose_file"
			}
			if !isConfigValid(conf) {
				os.Exit(int(syscall.EINVAL))
			}
			ctx := context.Background()
			code := proc.Run(ctx, conf, fileName)
			os.Exit(int(code))
		case "version":
			fmt.Println(version)
		default:
			fmt.Println("Unknown command, use help to view run examples")
			fmt.Println(usageInfo)
		}
	} else {
		fmt.Println("No command provided, use 'quack help' command to view usage information")
		fmt.Println(usageInfo)
	}
}

func ParseFlags() *utils.ConfigYaml {
	var conf utils.ConfigYaml

	notFound := conf.ReadConfig()
	if notFound != nil {
		fmt.Printf("Can't find config file quack_config.yaml: %s\n", notFound)
	}

	flag.Var(&conf.Models.Path, "models", "path to gorm models")
	flag.Var(&conf.Database.Uri, "uri", "database URI")
	flag.Var(&conf.Database.Name, "dbname", "database name")
	flag.Var(&conf.Migrations.Path, "path", "path to directory with migration files")

	flag.Var(&conf.Models.Exclude, "exclude", "Exclude gorm models")
	flag.Var(&conf.Database.Exclude, "db-exclude", "Exclude db tables")

	flag.Parse()

	if notFound != nil {
		if conf.Database.Uri == "" || conf.Database.Name == "" || (conf.Models.Path == "" && conf.Migrations.Path == "") {
			fmt.Println("Please provide all mandatory params(uri, models, dbname and path), using flags or configuration file")
			fmt.Println(usageInfo)
			os.Exit(1)
		}
	}
	return &conf
}

func isConfigValid(conf *utils.ConfigYaml) bool {
	// validate database Uri
	dbUriParts := strings.Split(conf.Database.Uri.String(), ":")
	if len(dbUriParts) == 1 {
		return false
	}
	conf.Database.Type = dbUriParts[0]
	return true
}

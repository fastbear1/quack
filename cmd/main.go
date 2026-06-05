package main

import (
	"flag"
	"fmt"

	utils "github.com/fastbear1/quack/internal/utils"
	proc  "github.com/fastbear1/quack/runner"
)

const (
	todoList = `TODO:
- parse config file
- parse command line args
- parse command line flags
- code structure
- two way compare state algorithm
- new lines in migration file
- create index statement
- create constraint statement
- parse all gorm tags
- Readme file
- tests for code function
- intergation tests with additional data files
- yaml config file parser
- parser for input args
- interface for all tables, column , index, constraints 
- independent drivers`

	helpInfo = `Quack - generate migration file for goose according gorm struct models 
information and database state. Use config file quack_config.yaml for running params.
Uasge:
  quack [flags] command
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
    - quack --modesl=internal/models --path=models --uri=postgres://user:pass@host:port/database --exclude=Base --db-exclude=goose,goose_migrations run`
)

func main() {
	var conf utils.ConfigYaml

	notFound := conf.ReadConfig()
	if notFound != nil {
		fmt.Println("Can't find config file quack_config.yaml")
	}

	flag.Var(&conf.Models.Path, "models", "path to gorm models")
	flag.Var(&conf.Database.Uri, "uri", "database URI")
	flag.Var(&conf.Database.Name, "dbname", "database name")
	flag.Var(&conf.Migrations.Path, "path", "path tp directory with migration files")

	flag.Var(&conf.Models.Exclude, "exclude", "Exlude gorm models")
	flag.Var(&conf.Database.Exclude, "db-exclude", "Exlude db tables")

	flag.Parse()

	if notFound != nil {
		if conf.Database.Uri == "" || conf.Database.Name == "" || conf.Models.Path == "" && conf.Migrations.Path == "" {
			fmt.Println("Please provide all mandatory params(uri, models, dbname and path), using flags or configuration file")
			panic("Exiting....")
		}
	}
	// check commands
	var commands []string = flag.Args()

	if len(commands) > 0 {
		switch commands[0] {
		case "help":
			fmt.Println(helpInfo)
		case "todo":
			fmt.Println(todoList)
		case "run":
			fmt.Println("Quacking migration file")
			proc.Run(&conf)
		default:
			fmt.Println("Unknown command, use help to view run exmaples")
		}
	} else {
		fmt.Println("Any command presented, use help to view usefull information")
	}

}

package main

import (
	"flag"
	"fmt"

	utils "github.com/fastbear1/quack/internal/utils"
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
	var yconf utils.ConfigYaml
	cnf, err := yconf.ReadConfig()
	if err != nil {
		fmt.Println("Can't config file")
	}

	flag.Var(&cnf.Models.Path, "models", "path to gorm models")
	flag.Var(&cnf.Database.Uri, "uri", "database URI")
	flag.Var(&cnf.Database.Name, "dbname", "database name")
	flag.Var(&cnf.Migrations.Path, "path", "path tp directory with migration files")

	flag.Var(&cnf.Models.Exclude, "exclude", "Exlude gorm models")
	flag.Var(&cnf.Database.Exclude, "db-exclude", "Exlude db tables")

	flag.Parse()
	fmt.Println(cnf)

	// check commands
	var commands []string = flag.Args()

	var conf utils.Config
	conf.GetConfig()
	fmt.Println(conf)

	if len(commands) > 0 {
		switch commands[0] {
		case "help":
			fmt.Println(helpInfo)
		case "todo":
			fmt.Println(todoList)
		case "run":
			fmt.Println("Quacking migration file")
		default:
			fmt.Println("Any command presented, use help to view usefull information")
		}
	} else {
		fmt.Println("Any command presented, use help to view usefull information")
	}

}

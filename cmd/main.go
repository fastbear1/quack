package main

import (
	"flag"
	"fmt"
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
	quack command [flags]
Commands:
	- run - quack(run and create) goose migration file
	- help - show help information
	- version - show current version(can be used for checking config file)
flags:
	- 
`
)

func main() {
	//todo := flag.String("h", "", "Print todo list")
	numbPtr := flag.Int("numb", 0, "int number")
	flag.Parse()
	if *numbPtr == 0 {
		fmt.Println("Numbet not defined")
	} else {
		fmt.Println(*numbPtr)
	}
	var commands []string = flag.Args()
	fmt.Println(commands)
	fmt.Println("tail:", commands)

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

package main

import (
	"flag"
	"fmt"
)

const todoList = `TODO:
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

func main() {
	//todo := flag.String("h", "", "Print todo list")
	numbPtr := flag.Int("numb", 0, "int number")
	flag.Parse()
	if *numbPtr == 0 {
		fmt.Println("Numbet not defined")
	} else {
		fmt.Println(*numbPtr)
	}
	fmt.Println(flag.Args())
	fmt.Println("tail:", flag.Args())

	for _, w := range flag.Args() {
		switch w {
		case "help":
			fmt.Println("Print help")
		case "todo":
			fmt.Println(todoList)
		}
	}
}

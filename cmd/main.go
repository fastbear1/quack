package main

import (
	"fmt"
)

func main() {
	var todoList = `
TODO:
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
	- independent drivers
`
	fmt.Println(todoList)
}

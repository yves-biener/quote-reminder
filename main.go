package main

import (
	"fmt"
	db "quote/db"
)

func main() {
	fmt.Println("Hello World")
	database := db.Connect("./test.sqlite")
	fmt.Printf("%v, %T\n", database, database)
}

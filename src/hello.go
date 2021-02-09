package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Name struct {
	Id   int
	Name string
}

func main() {
	db, err := sql.Open("mysql", "root:launder-motive-DREAR@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM names")

	if err != nil {
		log.Fatal(err)
	}

	var name Name
	err = row.Scan(&name.Id, &name.Name)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello, %v!\n", name.Name)
}

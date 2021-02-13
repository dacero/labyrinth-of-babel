package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Name struct {
	Id   int
	Name string
}

// initialises the database
func init() {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/")
	if err != nil {
		log.Fatal("Error when opening DB: ", err)
	}
	defer db.Close()

	file, err := ioutil.ReadFile("./hello.sql")
	if err != nil {
		// handle error
		log.Fatal("Error when initializing DB: ", err)
	}

	for _, request := range strings.Split(string(file), ";") {
		request = strings.TrimSpace(request)
		if request != "" {
			_, err := db.Exec(request)
			if err != nil {
				log.Fatalf("Error when initializing DB\n Query %s returned error: %s", request, err)
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/hello")
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

	fmt.Fprintf(w, "Hello, %v!\n", name.Name)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

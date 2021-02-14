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

type Cell struct {
	Id    string
	Title string
	Body  string
}

// initialises the database
func init() {
	log.Print("Initializing db... ")
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/")
	if err != nil {
		log.Fatal("Error when opening DB: ", err)
	}
	defer db.Close()

	file, err := ioutil.ReadFile("./db/labyrinth_of_babel.sql")
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
	log.Print("Finished initiatlizing db!")
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, title, body FROM cells")

	if err != nil {
		log.Fatal(err)
	}

	var cell Cell
	err = row.Scan(&cell.Id, &cell.Title, &cell.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "ID: %v!\n", cell.Id)
	fmt.Fprintf(w, "Title: %v!\n", cell.Title)
	fmt.Fprintf(w, "Body: %v!\n", cell.Body)

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

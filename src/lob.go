package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Cell struct {
	Id          string
	Title       string
	Body        string
	Create_time time.Time
	Update_time time.Time
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
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, title, body, create_time, update_time FROM cells")
	if err != nil {
		log.Fatal(err)
	}

	var cell Cell
	err = row.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Create_time, &cell.Update_time)
	if err != nil {
		log.Fatal(err)
	}

	t, _ := template.ParseFiles("./src/card.html")
	err = t.Execute(w, cell)
	if err != nil {
		log.Printf("Error when returning card: %s", err)
	}

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/dacero/labyrinth-of-babel/repository"

	_ "github.com/go-sql-driver/mysql"
)

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

func viewHandler(lob repository.LobRepository) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := r.URL.Path[len("/view/"):]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, _ := ioutil.ReadFile("./templates/card_not_found.html")
			fmt.Fprintf(w, string(notFound))
		} else {
			t, _ := template.ParseFiles("./templates/card.gohtml")
			err = t.Execute(w, cell)
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
		}
	})
}

func main() {
	lobRepository := repository.NewLobRepository()
	defer lobRepository.Close()
	http.HandleFunc("/view/", viewHandler(lobRepository))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

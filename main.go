package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/handlers"
	
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

func main() {
	lobRepository := repository.NewLobRepository()
	defer lobRepository.Close()
	http.HandleFunc("/view/", handlers.ViewHandler(lobRepository))
	http.HandleFunc("/new", handlers.CreateHandler(lobRepository))
	http.HandleFunc("/sources", handlers.SourcesHandler(lobRepository))
	http.HandleFunc("/rooms", handlers.RoomsHandler(lobRepository))
	http.HandleFunc("/page/", handlers.PageHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

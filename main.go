package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"net/http"
	"strings"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/handlers"
	"github.com/gorilla/mux"
	
	_ "github.com/go-sql-driver/mysql"
)

// initialises the database
func init() {
	log.Print("Initializing db... ")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	db, err := sql.Open("mysql", "root:"+password+"@tcp(mysql:3306)/")
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
	r := mux.NewRouter()
	r.HandleFunc("/cell/{id}", handlers.ViewHandler(lobRepository))
	r.HandleFunc("/cell/{id}/edit", handlers.EditHandler(lobRepository))
	r.HandleFunc("/cell/{id}/sources", handlers.SourcesHandler(lobRepository))
	r.HandleFunc("/cell/{id}/addSource", handlers.AddSourceHandler(lobRepository)).Methods("POST") //addSource
	r.HandleFunc("/cell/{id}/removeSource", handlers.RemoveSourceHandler(lobRepository)).Methods("POST") //removeSource
	r.HandleFunc("/cell/{id}/links", handlers.LinksHandler(lobRepository))
	r.HandleFunc("/cell/{id}/linkCell", handlers.LinkCellsHandler(lobRepository)).Methods("POST")
	r.HandleFunc("/cell/{id}/unlinkCell", handlers.UnlinkCellsHandler(lobRepository)).Methods("POST")
	r.HandleFunc("/save", handlers.SaveHandler(lobRepository))
	r.HandleFunc("/new", handlers.CreateHandler(lobRepository))
	r.HandleFunc("/searchSources", handlers.SearchSourcesHandler(lobRepository))
	r.HandleFunc("/searchRooms", handlers.SearchRoomsHandler(lobRepository))
	r.HandleFunc("/searchCells", handlers.SearchCellsHandler(lobRepository))
	r.HandleFunc("/page/{page}", handlers.PageHandler())
	log.Fatal(http.ListenAndServe(":8080", r))
}

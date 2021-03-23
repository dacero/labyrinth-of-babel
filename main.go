package main

import (
	"log"
	"net/http"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/handlers"
	"github.com/gorilla/mux"
	
	_ "github.com/go-sql-driver/mysql"
)

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
	r.HandleFunc("/rooms", handlers.RoomListHandler(lobRepository))
	log.Fatal(http.ListenAndServe(":8080", r))
}

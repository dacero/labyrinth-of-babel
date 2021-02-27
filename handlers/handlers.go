package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/dacero/labyrinth-of-babel/repository"
)

func ViewHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := r.URL.Path[len("/view/"):]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/card.gohtml")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			err = t.Execute(w, cell)
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
		}
	})
}

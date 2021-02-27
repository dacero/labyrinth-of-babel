package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/dacero/labyrinth-of-babel/repository"
)

func ViewHandler(lob repository.LobRepository) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := r.URL.Path[len("/view/"):]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, _ := ioutil.ReadFile("./templates/card_not_found.html")
			w.WriteHeader(http.StatusNotFound)
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

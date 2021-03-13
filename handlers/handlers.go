package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/models"
	"github.com/gorilla/mux"
)

func ViewHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
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

func EditHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/edit_card.gohtml")
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

func EditSourcesHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/edit_sources.gohtml")
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

func AddSourceHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		newSource := models.Source{Source: r.PostFormValue("source") }
		_, err := lob.AddSourceToCell(cellId, newSource)
		if err != nil {
			log.Printf("Error when adding source: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			http.Redirect(w, r, "/cell/"+cellId+"/edit/sources", http.StatusFound)
		}
	})
}

func RemoveSourceHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		source := models.Source{Source: r.PostFormValue("source") }
		_, err := lob.RemoveSourceFromCell(cellId, source)
		if err != nil {
			log.Printf("Error when removing source: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when removing source: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			http.Redirect(w, r, "/cell/"+cellId+"/edit/sources", http.StatusFound)
		}
	})
}

func SaveHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the form and create the cell
		updateCell := models.Cell{Id: r.PostFormValue("cellId"),
			Title: r.PostFormValue("title"),
			Body: r.PostFormValue("body"),
			Room: r.PostFormValue("room")}
		//call repository to create it
		_, err := lob.UpdateCell(updateCell)
		if err != nil {
			log.Printf("Error when updating card: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error when updating card: %s", err)
		} 
		//redirect to view the new cell card
		http.Redirect(w, r, "/cell/"+r.PostFormValue("cellId"), http.StatusFound)
	})
}

func CreateHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the form and create the cell
		log.Printf("New cell title: %s", r.PostFormValue("title"))
		newCell := models.Cell{Title: r.PostFormValue("title"),
			Body: r.PostFormValue("body"),
			Room: r.PostFormValue("room"),
			Sources: []models.Source{ models.Source{Source:r.PostFormValue("source")} } }
		//call repository to create it
		newCellId, err := lob.NewCell(newCell)
		if err != nil {
			log.Printf("Error when creating card: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error when creating card: %s", err)
		} 
		//redirect to view the new cell card
		http.Redirect(w, r, "/cell/"+newCellId+"/edit", http.StatusFound)
	})
}

func PageHandler() func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageName := mux.Vars(r)["page"]
		page, err := ioutil.ReadFile("./templates/" + pageName)
		if err != nil {
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning page: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(page))
	})
}

func SourcesHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.FormValue("term")
		sources := lob.SearchSources(term)
		returnString := "["
		for _, source := range sources {
			returnString += `"` + source.String() + `",`
		}
		returnString = returnString[:len(returnString)-1] + "]"
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, returnString)
	})
}

func RoomsHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.FormValue("term")
		rooms := lob.SearchRooms(term)
		returnString := "["
		for _, room := range rooms {
			returnString += `"` + room + `",`
		}
		returnString = returnString[:len(returnString)-1] + "]"
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, returnString)
	})
}
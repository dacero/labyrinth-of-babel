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

type Source struct {
	Source string
}

type CellLink struct {
	Id   string
	Text string
}

type Cell struct {
	Id          string
	Title       string
	Body        string
	Room        string
	Create_time time.Time
	Update_time time.Time
	Sources     []Source
	Links       []CellLink
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

func getCell(id string) Cell {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, title, body, room, create_time, update_time FROM cells WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	var cell Cell
	err = row.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Room, &cell.Create_time, &cell.Update_time)
	if err != nil {
		log.Fatal(err)
	}

	cell.Sources = getCellSources(id)
	cell.Links = getCellLinks(id)

	return cell
}

func getCellSources(id string) []Source {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT s.source 
		FROM sources s, cells_sources cs
		WHERE cs.sources_source = s.source 
		AND cs.cells_id=?`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var source string
		err := rows.Scan(&source)
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, Source{Source: source})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return sources
}

func getCellLinks(id string) []CellLink {
	db, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT c.id, c.title, c.body, c.update_time 
		FROM cells_links l, cells c 
		WHERE l.cells_a = c.id
		AND l.cells_b = ?
		UNION
		SELECT c.id, c.title, c.body, c.update_time 
		FROM cells_links l, cells c 
		WHERE l.cells_b = c.id
		AND l.cells_a = ?
		ORDER BY update_time DESC;`, id, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var links []CellLink
	for rows.Next() {
		var id, title, body, text string
		var update_time time.Time
		err := rows.Scan(&id, &title, &body, &update_time)
		if err != nil {
			log.Fatal(err)
		}
		if title == "" {
			r := []rune(body)
			if len(r) > 60 {
				text = string(r[0:50]) + "..."
			} else {
				text = body
			}
		} else {
			text = title
		}
		links = append(links, CellLink{Id: id, Text: text})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}

func handler(w http.ResponseWriter, r *http.Request) {

	cell := getCell("72aed05b-cb2d-4cad-bf70-05d8ae02a7bc")

	t, _ := template.ParseFiles("./src/card.gohtml")
	err := t.Execute(w, cell)
	if err != nil {
		log.Printf("Error when returning card: %s", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/dacero/labyrinth-of-babel/models"

	_ "github.com/go-sql-driver/mysql"
)

type LobRepository interface {
	GetCell(id string) (models.Cell, error)
	Close()
}

func NewLobRepository() *lobRepository {
	newDB, err := sql.Open("mysql", "root:secret@tcp(mysql:3306)/labyrinth_of_babel?parseTime=true")
	if err != nil {
		log.Panic(err)
	}
	return &lobRepository{db: newDB}
}

type lobRepository struct {
	db *sql.DB
}

func (r *lobRepository) getDB() *sql.DB {
	return r.db
}

func (r *lobRepository) Close() {
	r.db.Close()
}

func (r *lobRepository) GetCell(id string) (models.Cell, error) {
	var cell models.Cell

	row := r.getDB().QueryRow("SELECT id, title, body, room, create_time, update_time FROM cells WHERE id=?", id)

	err := row.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Room, &cell.Create_time, &cell.Update_time)
	if err != nil {
		log.Println(err)
		return cell, err
	}

	cell.Sources = r.getCellSources(id)
	cell.Links = r.getCellLinks(id)

	return cell, nil
}

func (r *lobRepository) getCellSources(id string) []models.Source {
	var sources []models.Source

	rows, err := r.getDB().Query(`SELECT s.source 
		FROM sources s, cells_sources cs
		WHERE cs.sources_source = s.source 
		AND cs.cells_id=?`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		err := rows.Scan(&source)
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, models.Source{Source: source})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return sources
}

func (r *lobRepository) getCellLinks(id string) []models.CellLink {
	var links []models.CellLink

	rows, err := r.getDB().Query(`SELECT c.id, c.title, c.body, c.update_time 
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
		links = append(links, models.CellLink{Id: id, Text: text})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}

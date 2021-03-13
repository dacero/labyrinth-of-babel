package repository

import (
	"database/sql"
	"log"
	"time"
	"errors"
	"strings"

	"github.com/dacero/labyrinth-of-babel/models"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

type LobRepository interface {
	//gets a new cell from its id
	GetCell(id string) (models.Cell, error)
	//updates the cell with new content
	UpdateCell(cell models.Cell) (int64, error)
	//adds and removes a source from a cell, returns the cell with updated sources
	AddSourceToCell(cellId string, source models.Source) (models.Cell, error)
	RemoveSourceFromCell(cellId string, source models.Source) (models.Cell, error)
	//creates a new cell and returns its new id
	NewCell(c models.Cell) (string, error)
	//searches for sources that contain the terms passed
	SearchSources(term string) []models.Source
	//searches for rooms that contain the terms passed
	SearchRooms(term string) []string
	//closes the database
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

func (r *lobRepository) UpdateCell(cell models.Cell) (int64, error) {
	//check the room and body to not be empty
	if strings.TrimSpace(cell.Room) == "" {
		return 0, errors.New("Empty room")
	}
	if strings.TrimSpace(cell.Body) == "" {
		return 0, errors.New("Empty body")
	}
	//insert the room first, just in case we need to create one
	err := r.insertRoom(cell.Room)
	if err != nil {
		return 0, err
	}
	//update the cell
	result, err := r.getDB().Exec("UPDATE cells SET title = ?, body = ?, room = ?, update_time = ? where id = ?", cell.Title, cell.Body, cell.Room, time.Now(), cell.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *lobRepository) AddSourceToCell(cellId string, source models.Source) (models.Cell, error) {
	//create the source if not exists
	err := r.insertSources([]models.Source{source})
	if err != nil {
		cell, _ := r.GetCell(cellId) // unnecessary!!!
		return cell, err
	}
	//link the source to the cell if not already
	err = r.linkSources(cellId, []models.Source{source})
	if err != nil {
		cell, _ := r.GetCell(cellId) // unnecessary!!!
		return cell, err
	}
	return r.GetCell(cellId)
}

func (r *lobRepository) RemoveSourceFromCell(cellId string, source models.Source) (models.Cell, error) {
	stmt, err := r.getDB().Prepare("DELETE FROM cells_sources WHERE cells_id=? AND sources_source=?")
	if err != nil {
		cell, _ := r.GetCell(cellId) // unnecessary!!!
		return cell, err
	}
	_, err = stmt.Exec(cellId, strings.TrimSpace(source.Source))
	if err != nil {
		cell, _ := r.GetCell(cellId) // unnecessary!!!
		return cell, err
	}
	return r.GetCell(cellId)
}

func (r *lobRepository) NewCell(cell models.Cell) (string, error) {
	//validations
	if strings.TrimSpace(cell.Room) == "" {
		return "", errors.New("Empty room name")
	}
	if strings.TrimSpace(cell.Body) == "" {
		return "", errors.New("Empty body in new cell")
	}
	
	cellId := uuid.NewString()
	
	//insert the room
	err := r.insertRoom(cell.Room)
	if err != nil {
		return "", err
	}
	//insert the cell
	stmt, err := r.getDB().Prepare("INSERT INTO cells(id, title, body, room) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(cellId, strings.TrimSpace(cell.Title), strings.TrimSpace(cell.Body), strings.TrimSpace(cell.Room))
	if err != nil {
		return "", err
	}
	//insert the sources
	err = r.insertSources(cell.Sources)
	if err != nil {
		return "", err
	}
	//link sources with the cell
	err = r.linkSources(cellId, cell.Sources)
	if err != nil {
		return "", err
	}
	return cellId, nil
}

func (r *lobRepository) insertSources(sources []models.Source) error {
	insertStr := "INSERT IGNORE INTO sources(source) VALUES "
	vals := []interface{}{}
	for _, source := range sources {
		insertStr += "(?),"
		vals = append(vals, strings.TrimSpace(string(source.String())))
	}
	//trim the last
	insertStr = insertStr[:len(insertStr)-1]
	stmt, err := r.getDB().Prepare(insertStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err
	}
	return nil	
}

func (r *lobRepository) insertRoom(room string) error {
	if strings.TrimSpace(room) == "" {
		return errors.New("Empty room name")
	}
	stmt, err := r.getDB().Prepare("INSERT IGNORE INTO rooms(room) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(strings.TrimSpace(room))
	if err != nil {
		return err
	}
	return nil	
}

func (r *lobRepository) linkSources(cellId string, sources []models.Source) error {
	insertStr := "INSERT IGNORE INTO cells_sources(cells_id, sources_source) VALUES "
	vals := []interface{}{}
	for _, source := range sources {
		insertStr += "(?, ?),"
		vals = append(vals, cellId, string(source.String()))
	}
	//trim the last
	insertStr = insertStr[:len(insertStr)-1]
	stmt, err := r.getDB().Prepare(insertStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err
	}
	return nil	
}

func (r *lobRepository) SearchSources(term string) []models.Source {
	var sources []models.Source
	
	rows, err := r.getDB().Query(`SELECT source 
		FROM sources
		WHERE source LIKE ?`, "%" + term + "%")
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

func (r *lobRepository) SearchRooms(term string) []string {
	var rooms []string
	
	rows, err := r.getDB().Query(`SELECT room 
		FROM rooms
		WHERE room LIKE ?`, "%" + term + "%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var room string
		err := rows.Scan(&room)
		if err != nil {
			log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rooms
}

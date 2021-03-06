package repository

import (
	"database/sql"
	"log"
	"os"
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
	//links or unlink 2 cells
	LinkCells(idA string, idB string) (error)
	UnlinkCells(idA string, idB string) (error)
	//checks if two cells are linked
	CheckLink(idA string, idB string) (bool, error)
	//creates a new cell and returns its new id
	NewCell(c models.Cell) (string, error)
	//Returns a full list of all rooms in the labyrinth
	ListRooms() ([]models.CollectionOfCells, error)
	//Returns all cells in a room
	ListCellsInRoom(room string) ([]models.Cell, error)
	//searches for sources that contain the terms passed
	SearchSources(term string) []models.Source
	//searches for rooms that contain the terms passed
	SearchRooms(term string) []string
	//searches for cells that contain the terms passed
	SearchCells(term string) []models.Cell
	//closes the database
	Close()
}

func NewLobRepository() *lobRepository {
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	newDB, err := sql.Open("mysql", "root:"+password+"@tcp(mysql:3306)/"+database+"?parseTime=true")
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

func (r *lobRepository) getCellLinks(id string) []models.Cell {
	var links []models.Cell

	rows, err := r.getDB().Query(`SELECT c.id, c.title, c.body, c.create_time, c.update_time, c.room 
		FROM cells_links l, cells c 
		WHERE l.cells_a = c.id
		AND l.cells_b = ?
		UNION
		SELECT c.id, c.title, c.body, c.create_time, c.update_time, c.room
		FROM cells_links l, cells c 
		WHERE l.cells_b = c.id
		AND l.cells_a = ?
		ORDER BY create_time DESC;`, id, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cell models.Cell
		err := rows.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Create_time, &cell.Update_time, &cell.Room)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, cell)
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

func (r *lobRepository) LinkCells(idA string, idB string) (error) {
	//verify that the cells are not already linked
	if idA == idB {
		return errors.New("Tried linking a cell with itself")
	}
	linked, err := r.CheckLink(idA, idB)
	if err != nil {
		return err
	}
	if linked {
		return errors.New("Tried linking cells already linked")
	}
	//link the cells
	stmt, err := r.getDB().Prepare("INSERT INTO cells_links(cells_a, cells_b) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(idA, idB)
	return err
}

func (r *lobRepository) UnlinkCells(idA string, idB string) (error) {
	//link the cells
	stmt, err := r.getDB().Prepare("DELETE FROM cells_links WHERE (cells_a = ? AND cells_b = ?) OR (cells_a = ? AND cells_b = ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(idA, idB, idB, idA)
	return err
}

func (r *lobRepository) CheckLink(idA string, idB string) (bool, error) {
	var numLinks int
	row := r.getDB().QueryRow("SELECT COUNT(*) FROM cells_links WHERE (cells_a=? AND cells_b=?) OR (cells_a=? AND cells_b=?)", idA, idB, idB, idA)
	err := row.Scan(&numLinks)
	if err != nil {
		log.Print(err)
		return false, err
	}
	if numLinks == 0 { return false, nil }
	if numLinks == 1 { return true, nil }
	return true, errors.New("Too many links")
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
		if trimmedSource := strings.TrimSpace(string(source.String())); trimmedSource != "" {
			insertStr += "(?),"
			vals = append(vals, trimmedSource)
		}
	}
	if len(vals) == 0 { return nil }
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
		if trimmedSource := strings.TrimSpace(string(source.String())); trimmedSource != "" {
			insertStr += "(?, ?),"
			vals = append(vals, cellId, trimmedSource)
		}
	}
	if len(vals) == 0 { return nil }
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

func (r *lobRepository) SearchCells(term string) []models.Cell {
	var cells []models.Cell
	rows, err := r.getDB().Query(`SELECT id, title, body, create_time, update_time, room
		FROM cells
		WHERE title LIKE ? OR body LIKE ?`, "%" + term + "%", "%" + term + "%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cell models.Cell
		err := rows.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Create_time, &cell.Update_time, &cell.Room)
		if err != nil {
			log.Fatal(err)
		}
		cells = append(cells, cell)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return cells
}

func (r *lobRepository) ListRooms() ([]models.CollectionOfCells, error) {
	var rooms []models.CollectionOfCells
	
	rows, err := r.getDB().Query(`SELECT rooms.room, COUNT(*), MIN(create_time) create_time 
		FROM rooms, cells
		WHERE rooms.room = cells.room
		GROUP BY rooms.room
		ORDER BY create_time DESC`)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.CollectionOfCells
		err := rows.Scan(&room.Name, &room.CellCount, &room.Create_time)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (r *lobRepository) ListCellsInRoom(room string) ([]models.Cell, error) {
	var cells []models.Cell
	
	rows, err := r.getDB().Query(`SELECT id, title, body, room, create_time, update_time 
		FROM cells 
		WHERE room=?
		ORDER BY update_time DESC`, room)
	if err != nil {
		return cells, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var cell models.Cell
		err := rows.Scan(&cell.Id, &cell.Title, &cell.Body, &cell.Room, &cell.Create_time, &cell.Update_time)
		if err != nil {
			return cells, err
		}
		cells = append(cells, cell)
	}
	if err := rows.Err(); err != nil {
		return cells, err
	}
	
	return cells, nil
}	

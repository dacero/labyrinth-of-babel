package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Name struct {
	Id   int
	Name string
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:launder-motive-DREAR@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM names")

	if err != nil {
		log.Fatal(err)
	}

	var name Name
	err = row.Scan(&name.Id, &name.Name)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Hello, %v!\n", name.Name)
}

func main() {
	// first we ensure we have a fresh db install
	db, err := sql.Open("mysql", "root:launder-motive-DREAR@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS hello")
	if err != nil {
		log.Printf("Error %s when deleting DB\n", err)
		return
	}
	_, err = db.ExecContext(ctx, "CREATE DATABASE hello")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	_, err = db.ExecContext(ctx, "USE hello")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	_, err = db.ExecContext(ctx, `CREATE TABLE names (
		  id int unsigned NOT NULL AUTO_INCREMENT,
		  name varchar(255) NOT NULL DEFAULT '',
		  PRIMARY KEY (id))`)
	if err != nil {
		log.Printf("Error %s when creating the table\n", err)
		return
	}
	_, err = db.ExecContext(ctx, "INSERT INTO names VALUES (1, 'Created on the fly!')")
	if err != nil {
		log.Printf("Error %s when inserting in the table\n", err)
		return
	}
	db.Close()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

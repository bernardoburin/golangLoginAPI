package database

import(
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"time"
)

func CreateDatabaseConnection() *sql.DB {
	connStr := "user=postgres password=postgres dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected to the database successfully!")
	return db
}
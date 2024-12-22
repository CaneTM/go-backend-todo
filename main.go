package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/canetm/go-backend-todo/common"
	"github.com/canetm/go-backend-todo/handlers"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db := connectDB()
	addHandlers(db)
	fmt.Printf("Server listening on http://localhost%s\n", common.Port)
	http.ListenAndServe(common.Port, nil)
}

func connectDB() *sql.DB {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "GoBackendTodo",
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Printf("Connected to %s database\n", cfg.DBName)
	return db
}

func addHandlers(db *sql.DB) {
	handlers.NewUserHandler(db).HandleService()
}


//Robert was here 
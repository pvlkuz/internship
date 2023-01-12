package main

import (
	"log"
	"main/crud_handler"
	database "main/data-base"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	m, err := migrate.New("file://./migration", "postgresql://postgres:password@database:5432/postgres?sslmode=disable")
	if err != nil {
		time.Sleep(2 * time.Second)
		m, err = migrate.New("file://./migration", "postgresql://postgres:password@database:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatalf("failed to migration init: %s", err.Error())
		}
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("failed to migrate up: %s", err.Error())
	}

	db, err := database.NewDB()
	if err != nil {
		time.Sleep(2 * time.Second)
		db, err = database.NewDB()
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err.Error())
		}

	}

	handler := crud_handler.NewHandler(*db)
	handler.RunServer()

}

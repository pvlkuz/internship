package main

import (
	"log"
	database "main/DataBase"
	"main/crud_handler"
)

func main() {

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	handler := crud_handler.NewHandler(*db)
	handler.RunServer()

}

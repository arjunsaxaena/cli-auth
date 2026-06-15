package main

import (
	"fmt"
	"log"

	"cli-auth/internal/database"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = database.CreateSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully.")
}

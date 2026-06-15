package main

import (
	"cli-auth/internal/auth"
	"cli-auth/internal/cli"
	"cli-auth/internal/database"
	"log"
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

	repo := auth.NewRepository(db)

	err = cli.Start(repo)
	if err != nil {
		log.Fatal(err)
	}
}

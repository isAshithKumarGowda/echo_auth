package main

import (
	"database/sql"
	"log"

	"github.com/isAshithKumarGowda/Echo_Auth/pkg/database"
)

func Start(db *sql.DB, addr *string) {
	e := InitialiseHttpRouter(db)

	query := database.NewDBinstance(db)
	err := query.InitialiseDBqueries()
	if err != nil {
		log.Fatalf("Unable to Initialize Database :8080")
	}
	log.Printf("Starting server at %s", *addr)
	e.Start(*addr)
}

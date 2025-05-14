package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Connection struct {
	db *sql.DB
}

func NewDatabase() *Connection {
	DB, err := sql.Open("postgres", "postgresql://snipetbox_user:Sbqx0etnAvHSIvCBDhCBfda6G5iXwdVT@dpg-d0i1kd24d50c73b0l9b0-a.singapore-postgres.render.com/snipetbox")
	if err != nil {
		panic(err)
	}
	return &Connection{
		db: DB,
	}
}

func (c *Connection) CheckStatus() {
	err := c.db.Ping()
	if err != nil {
		log.Fatalf("Bad Database connection %v", err)
	}
	log.Println("Connected to database successfully")
}

func (c *Connection) Close() {
	if err := c.db.Close(); err != nil {
		log.Fatalf("Error closing the database")
	}
	log.Println("Database closed")
}

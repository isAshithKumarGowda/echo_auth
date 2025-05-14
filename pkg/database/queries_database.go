package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Query struct {
	db *sql.DB
}

func NewDBinstance(db *sql.DB) *Query {
	return &Query{
		db: db,
	}
}

func (db *Query) InitialiseDBqueries() error {

	queries := []string{
		`CREATE TABLE IF NOT EXISTS admin (
			admin_id VARCHAR(36) PRIMARY KEY,
			admin_name VARCHAR(100) NOT NULL,
			admin_email VARCHAR(100) NOT NULL UNIQUE,
			admin_hash_password VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(36) PRIMARY KEY,
			user_name VARCHAR(100) NOT NULL,
			user_email VARCHAR(100) NOT NULL UNIQUE,
			user_hash_password VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			email VARCHAR(100) UNIQUE NOT NULL,
			otp VARCHAR(10) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS verified_emails (
			id SERIAL PRIMARY KEY,
			email VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS admin_login_history (
			email VARCHAR(100) PRIMARY KEY,
			login_time TIMESTAMPTZS DEFAULT NOW(),
			logout_time TIMESTAMPTZS
		)`,
		`CREATE TABLE IF NOT EXISTS user_login_history (
			email VARCHAR(100) PRIMARY KEY,
			login_time TIMESTAMPTZS DEFAULT NOW(),
			logout_time TIMESTAMPTZS
		)`,
	}

	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("error while initialising DB %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			log.Println("Initialised Database")
		}
	}()

	for _, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil

}

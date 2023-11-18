package server

import (
	"os"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db *sqlx.DB
}

func newApp() *App {
	dbURL := os.Getenv("PRODUCTAPI_DB_URL")
	if dbURL == "" {
		panic("db url not specified")
	}
	db, err := connectPostgres(dbURL)
	if err != nil {
		panic("connect postgres error: " + err.Error())
	}

	return &App{
		db: db,
	}
}

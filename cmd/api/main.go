package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/3n0ugh/simple-crud-api/internal/data"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	dsn  string
}

type application struct {
	model data.Model
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.dsn, "db-dsn",
		"postgres://postgres@localhost/book?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	app := &application{
		data.NewModel(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.SetRouter(),
	}

	log.Fatal(srv.ListenAndServe())
}

func openDB(cfg config) (*sql.DB, error) {

	// Create an empty connection pool
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the connection couldn't be established successfully
	// within the 5-second deadline, then this will return an error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

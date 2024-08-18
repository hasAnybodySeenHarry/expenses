package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func openDB(cfg *db, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConn)
	db.SetMaxIdleConns(cfg.maxIdleConn)

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err = db.PingContext(ctx)
		cancel()

		if nil == err {
			return db, nil
		}

		log.Printf("Failed to connect to the database: %v. Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	return nil, err
}

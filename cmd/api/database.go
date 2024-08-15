package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func openDB(cfg *db) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConn)
	db.SetMaxIdleConns(cfg.maxIdleConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

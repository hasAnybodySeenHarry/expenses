package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

func initDependencies(cfg config, logger *log.Logger) (*sql.DB, *amqp.Connection, error) {
	var db *sql.DB
	var conn *amqp.Connection
	var dbErr, amqpErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		db, dbErr = openDB(&cfg.db, 10)
		if dbErr != nil {
			logger.Printf("Failed to connect to the database: %v", dbErr)
		} else {
			logger.Println("Successfully connected to the database")
		}
	}()

	go func() {
		defer wg.Done()
		conn, amqpErr = openAMQP(cfg.amqp, 6)
		if amqpErr != nil {
			logger.Printf("Failed to connect to the messaging proxy: %v", amqpErr)
		} else {
			logger.Println("Successfully connected to the messaging proxy")
		}
	}()

	wg.Wait()

	if dbErr != nil {
		return nil, nil, dbErr
	}
	if amqpErr != nil {
		return nil, nil, amqpErr
	}

	return db, conn, nil
}

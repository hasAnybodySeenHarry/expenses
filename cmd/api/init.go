package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/streadway/amqp"
)

func initDependencies(cfg config, logger *log.Logger) (*sql.DB, *amqp.Connection, sarama.SyncProducer, error) {
	var dbErr, amqpErr, producerErr error
	var db *sql.DB
	var conn *amqp.Connection
	var producer sarama.SyncProducer

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		db, dbErr = openDB(&cfg.db, 10)
		if dbErr != nil {
			logger.Printf("Failed to connect to the database: %v", dbErr)
		} else {
			logger.Println("Successfully connected to the database")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, amqpErr = openAMQP(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.msgProxy.username, cfg.msgProxy.password, cfg.msgProxy.host, cfg.msgProxy.port), 6)
		if amqpErr != nil {
			logger.Printf("Failed to connect to the messaging proxy: %v", amqpErr)
		} else {
			logger.Println("Successfully connected to the messaging proxy")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer, producerErr = openProducer([]string{fmt.Sprintf("%s:%d", cfg.pub.host, cfg.pub.port)}, 5)
		if producerErr != nil {
			logger.Printf("Failed to connect to the Kafka producer: %v", producerErr)
		} else {
			logger.Println("Successfully connected to the Kafka producer")
		}
	}()

	wg.Wait()

	if dbErr != nil {
		return nil, nil, nil, dbErr
	}
	if amqpErr != nil {
		return nil, nil, nil, amqpErr
	}
	if producerErr != nil {
		return nil, nil, nil, producerErr
	}

	return db, conn, producer, nil
}

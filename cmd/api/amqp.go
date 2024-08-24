package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func openAMQP(uri string, maxRetries int) (*amqp.Connection, error) {
	var err error
	var conn *amqp.Connection

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(uri)
		if err == nil {
			return conn, nil
		}
		log.Printf("Failed to connect to RabbitMQ: %s. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}

package main

import (
	"github.com/streadway/amqp"
)

func openAMQP(uri string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

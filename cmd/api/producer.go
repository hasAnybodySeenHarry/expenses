package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

func openProducer(brokers []string, retries int) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	var err error
	var producer sarama.SyncProducer

	for i := 0; i < retries; i++ {
		producer, err = sarama.NewSyncProducer(brokers, config)
		if nil == err {
			return producer, nil
		}

		log.Printf("Failed to connect to Kafka: %s. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}

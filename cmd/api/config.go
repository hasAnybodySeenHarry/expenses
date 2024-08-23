package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	port     int
	env      string
	db       db
	msgProxy msgProxy
	grpcPort int
	pub      publisher
}

type msgProxy struct {
	username string
	password string
	port     int
	host     string
}

type db struct {
	host        string
	name        string
	username    string
	password    string
	port        int
	maxOpenConn int
	maxIdleConn int
}

type publisher struct {
	host string
	port int
}

func loadConfig(cfg *config) {
	flag.IntVar(&cfg.grpcPort, "tcp", getEnvInt("GRPC_PORT", 50051), "The port that the grpc server listens at")

	flag.IntVar(&cfg.port, "port", getEnvInt("PORT", 4000), "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "The environment of the server")

	flag.StringVar(&cfg.db.host, "db-host", os.Getenv("DB_HOST"), "The address to connect to postgres")
	flag.StringVar(&cfg.db.name, "db-name", os.Getenv("DB_NAME"), "The name of the postgres database to connect to")
	flag.StringVar(&cfg.db.username, "db-username", os.Getenv("DB_USERNAME"), "The username to connect to postgres")
	flag.StringVar(&cfg.db.password, "db-password", os.Getenv("DB_PASSWORD"), "The password to connect to postgres")
	flag.IntVar(&cfg.db.port, "db-port", getEnvInt("DB_PORT", 5432), "The port to connect to postgres")
	flag.IntVar(&cfg.db.maxOpenConn, "max-open-conn", getEnvInt("MAX_OPEN_CONN", 30), "The maximum number of opened connections")
	flag.IntVar(&cfg.db.maxIdleConn, "max-idle-conn", getEnvInt("MAX_IDLE_CONN", 30), "The maximum number of idle connections")

	flag.StringVar(&cfg.msgProxy.username, "amqp-username", os.Getenv("AMQP_USERNAME"), "The username to connect to AMQP messaging proxy")
	flag.StringVar(&cfg.msgProxy.password, "amqp-password", os.Getenv("AMQP_PASSWORD"), "The password to connect to AMQP messaging proxy")
	flag.StringVar(&cfg.msgProxy.host, "amqp-host", os.Getenv("AMQP_HOST"), "The address to connect to AMQP messaging proxy")
	flag.IntVar(&cfg.msgProxy.port, "amqp-port", getEnvInt("AMQP_PORT", 5672), "The port to connect to AMQP messaging proxy")

	flag.StringVar(&cfg.pub.host, "pub-host", os.Getenv("PUB_HOST"), "The address to connect to Kafka node")
	flag.IntVar(&cfg.pub.port, "pub-port", getEnvInt("PUB_PORT", 9092), "The port to connect to Kafka node")

	flag.Parse()
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Invalid value for environment variable %s: %s\n", key, valueStr)
		return defaultValue
	}

	return value
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	port int
	env  string
	db   db
	amqp string
}

type db struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
}

func loadConfig(cfg *config) {
	flag.IntVar(&cfg.port, "port", getEnvInt("PORT", 4000), "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "The environment of the server")

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "The datasource to connect to postgres")
	flag.IntVar(&cfg.db.maxOpenConn, "max-open-conn", getEnvInt("MAX-OPEN-CONN", 30), "The maximum number of opened connections")
	flag.IntVar(&cfg.db.maxIdleConn, "max-idle-conn", getEnvInt("MAX-IDLE-CONN", 30), "The maximum number of idle connections")

	flag.StringVar(&cfg.amqp, "amqp-uri", os.Getenv("AMQP_URI"), "The URI of AMQP messaging proxy")

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

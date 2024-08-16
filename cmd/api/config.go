package main

import (
	"flag"
	"os"
)

type config struct {
	port int
	env  string
	db   db
}

type db struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
}

func loadConfig(cfg *config) {
	// public.ecr.aws/docker/library/postgres:latest
	// postgres:alpine
	// postgres:latest

	// docker run --name my-postgres -d -p 5432:5432 -e POSTGRES_USER=harry -e POSTGRES_PASSWORD=password -e POSTGRES_DB=expenses

	flag.IntVar(&cfg.port, "port", 4000, "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", "development", "The environment of the server")

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("dsn"), "The datasource to connect to postgres")
	flag.IntVar(&cfg.db.maxOpenConn, "max-open-conn", 30, "The maximum number of opened connections")
	flag.IntVar(&cfg.db.maxIdleConn, "max-idle-conn", 30, "The maximum number of idle connections")

	flag.Parse()
}

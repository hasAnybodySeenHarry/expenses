package main

import "flag"

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
	flag.IntVar(&cfg.port, "port", 4000, "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", "development", "The environment of the server")

	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://harry:password@localhost/expenses?sslmode=disable", "The datasource to connect to postgres")
	flag.IntVar(&cfg.db.maxOpenConn, "max-open-conn", 30, "The maximum number of opened connections")
	flag.IntVar(&cfg.db.maxIdleConn, "max-idle-conn", 30, "The maximum number of idle connections")

	flag.Parse()
}

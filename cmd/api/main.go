package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/mailer"
)

type application struct {
	config config
	wg     sync.WaitGroup
	logger *log.Logger
	models data.Models
	mailer *mailer.Mailer
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(&cfg.db)
	if err != nil {
		l.Fatalln(err)
	}
	l.Println("Successfully connected to the database")
	defer db.Close()

	conn, err := openAMQP(cfg.amqp)
	if err != nil {
		log.Fatalln(err)
	}
	l.Println("Successfully connected to the messaging proxy")

	mailer, err := mailer.New(conn, "email_queue")
	if err != nil {
		l.Fatalln(err)
	}
	defer mailer.Close()

	app := application{
		config: cfg,
		logger: l,
		models: data.New(db),
		mailer: mailer,
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatalln(err)
	}
}

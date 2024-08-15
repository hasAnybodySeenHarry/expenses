package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/expenses/internal/data"
)

type application struct {
	config config
	wg     sync.WaitGroup
	logger *log.Logger
	models data.Models
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

	app := application{
		config: cfg,
		logger: l,
		models: data.New(db),
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatalln(err)
	}

}

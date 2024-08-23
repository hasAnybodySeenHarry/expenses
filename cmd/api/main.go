package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/mailer"
	"harry2an.com/expenses/internal/notifier"
)

type application struct {
	config    config
	wg        sync.WaitGroup
	logger    *log.Logger
	models    data.Models
	notifiers *notifier.Notifiers
	mailer    *mailer.Mailer
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, conn, producer, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer db.Close()
	defer conn.Close()
	defer producer.Close()

	mailer, err := mailer.New(conn, "email_queue")
	if err != nil {
		l.Fatalln(err)
	}
	defer mailer.Close()

	app := application{
		config:    cfg,
		logger:    l,
		models:    data.New(db),
		mailer:    mailer,
		notifiers: notifier.New(producer),
	}

	var servers sync.WaitGroup
	servers.Add(2)

	go func() {
		defer servers.Done()
		if err := app.grpc(cfg.grpcPort); err != nil {
			app.logger.Fatalln("gRPC server stopped with error:", err)
		}
	}()

	go func() {
		defer servers.Done()
		if err := app.serve(); err != nil {
			app.logger.Fatalln("HTTP server stopped with error:", err)
		}
	}()

	servers.Wait()
	app.logger.Println("Both servers have stopped gracefully.")
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      app.routes(),
	}

	app.logger.Println("Server is starting")

	shutdownErr := make(chan error, 1)
	go func() {
		relay := make(chan os.Signal, 1)
		signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		s := <-relay
		app.logger.Println("Received signal:", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		app.logger.Println("Starting to clean up goroutines")

		app.wg.Wait()
		shutdownErr <- err
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	app.logger.Println("The server has just stopped now")

	err = <-shutdownErr
	if err != nil {
		app.logger.Println(err)
		return err
	}

	app.logger.Println("The server has completely stopped")

	return nil
}

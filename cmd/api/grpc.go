package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (app *application) grpc(port int) error {
	server := grpc.NewServer()
	reflection.Register(server)

	app.registerGRPCservers(server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	relay := make(chan os.Signal, 1)
	signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM)

	shutdownErr := make(chan error)

	go func() {
		<-relay
		app.logger.Println("Received shutdown signal. Initiating graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.GracefulStop()

		select {
		case <-ctx.Done():
			app.logger.Println("Server graceful shutdown timed out")
			shutdownErr <- fmt.Errorf("error shutdown timed out")
		default:
			app.logger.Println("Server gracefully stopped")
			shutdownErr <- nil
		}
	}()

	app.logger.Printf("gRPC server listening on :%d", port)
	if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("failed to serve: %v", err)
	}

	err = <-shutdownErr
	return err
}

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	httpInfra "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"

	grcpInfra "ride-sharing/services/trip-service/internal/infrastructure/grpc"

	"google.golang.org/grpc"
)

const grpcAddr = ":9093"

func main() {

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	mux := http.NewServeMux()

	httpHandler := httpInfra.HttpHandler{
		Service: svc,
	}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	grpcServer := grpc.NewServer()

	tcpListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		return
	}

	defer tcpListener.Close()

	grcpInfra.NewGrpcHandler(grpcServer, svc)

	serverError := make(chan error, 2)

	go func() {
		fmt.Println("Server started on", server.Addr)

		serverError <- server.ListenAndServe()
	}()

	go func() {
		fmt.Println("gRPC server started on", grpcAddr)

		serverError <- grpcServer.Serve(tcpListener)
	}()

	ctxSignal, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	select {
	case err := <-serverError:
		fmt.Printf("Error starting server: %v\n", err)

	case <-ctxSignal.Done():
		fmt.Println("Received shutdown signal")
	}

	grpcDone := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcDone)
	}()

	select {
	case <-grpcDone:
		fmt.Println("gRPC server stopped")
	case <-time.After(5 * time.Second):
		fmt.Println("gRPC server did not stop in time, forcing exit")
		grpcServer.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Could not shut down server gracefully: %v\n", err)
		server.Close()
	}

}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	mux := http.NewServeMux()

	httpHandler := h.HttpHandler{
		Service: svc,
	}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	serverError := make(chan error, 1)

	go func() {
		fmt.Printf("Server started on %s\n", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		fmt.Printf("Error starting server: %v\n", err)
	case <-shutdown:
		fmt.Println("Server shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Could not shut down server gracefully: %v\n", err)
			server.Close()
		}
	}

	// if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("server failed to start: %v", err)
	// }
}

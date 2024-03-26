package main

import (
	"fmt"
	"internal/http-server/handlers/cards"
	"internal/model"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	model, err := model.New("http://127.0.0.1:5000/rec")
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	// TODO : add config
	router := chi.NewRouter()
	router.Get("/cards", cards.New(model))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "localhost:8085",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
		IdleTimeout:  40 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	<-done
}

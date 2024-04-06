package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_backend_frontend/internal/http-server/handlers/cards"
	sessions_handler "test_backend_frontend/internal/http-server/handlers/session"
	"test_backend_frontend/internal/model"
	sessions "test_backend_frontend/internal/sessions"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	model, err := model.New("http://0.0.0.0:5000/rec")
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	var sessionManager *sessions.SessionManager
	sessionManager, err = sessions.NewSessionManager("localhost:6379", "", 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// TODO : add config
	router := chi.NewRouter()
	router.Get("/cards", cards.New(model))
	router.Post("/sessions", sessions_handler.SessionCreatePage(sessionManager))
	router.Post("/sessions/{id}", sessions_handler.SessionGetData(sessionManager))
	router.Patch("/sessions/{id}", sessions_handler.SessionAdduser(sessionManager))
	router.Put("/sessions/{id}", sessions_handler.SessionModifyuser(sessionManager))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
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

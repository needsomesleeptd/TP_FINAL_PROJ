package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_backend_frontend/internal/http-server/handlers/cards"
	"test_backend_frontend/internal/model"
	"test_backend_frontend/internal/models"
	sessions "test_backend_frontend/internal/sessions"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TempSessionData struct {
	Session_id   uuid.UUID
	People_count int
}

var sessionManager *sessions.SessionManager

func session_create_page(w http.ResponseWriter, r *http.Request) {
	req := r.FormValue("request")
	userReq := models.NewUserReq(2, "anyname", req)
	sessionID, err := sessionManager.CreateSession(userReq)
	if err != nil {
		fmt.Errorf(err.Err().Error())
	}

	var tpl_index = template.Must(template.ParseFiles("../frontend/create_session.html"))

	tpl_index.Execute(w, nil)
	s_path := fmt.Sprintf("/session/%u", sessionID)
	http.Redirect(w, r, s_path, http.StatusFound)
}

func session_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sessionID uuid.UUID
	var err error
	sessionID, err = uuid.Parse(vars["id"])
	if err != nil {
		fmt.Print(err)
	}
	users, err := sessionManager.GetUsers(sessionID)
	if err != nil {
		fmt.Print(err)
	}
	var tpl_index = template.Must(template.ParseFiles("../frontend/sessions.html"))
	data := TempSessionData{Session_id: sessionID, People_count: len(users)}
	err = tpl_index.Execute(w, data)
	if err != nil {
		fmt.Print(err)
	}
}

func main() {
	model, err := model.New("http://0.0.0.0:5000/rec")
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	// TODO : add config
	router := chi.NewRouter()
	router.Get("/cards", cards.New(model))
	router.Post("/session_create", session_create_page)
	router.Get("/session/{id}", session_page)

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

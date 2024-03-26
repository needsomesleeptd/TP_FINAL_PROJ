package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_backend_frontend/internal/http-server/handlers/cards"
	"test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/model"
	"test_backend_frontend/internal/models"
	sessions "test_backend_frontend/internal/sessions"
	"time"

	resp "test_backend_frontend/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ResponseSessionID struct {
	Response  resp.Response
	SessionID uuid.UUID `json:"sessionID"`
}

type ResponseUsersReq struct {
	Response  resp.Response
	UsersReqs []models.UserReq
}

type RequestSessionUsers struct {
	SessionID uuid.UUID `json:"sessionID"`
}

type RequestAddUser struct {
	User      models.UserReq `json:"user"`
	SessionID uuid.UUID      `json:"sessionID"`
}

var sessionManager *sessions.SessionManager

func session_create_page(w http.ResponseWriter, r *http.Request) {

	userReq := models.NewUserReq(2, "anyname", "initializer of  a party")
	sessionID, err := sessionManager.CreateSession(userReq)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	render.JSON(w, r, ResponseSessionID{
		Response:  resp.OK(),
		SessionID: sessionID,
	})

}

func session_get_data(w http.ResponseWriter, r *http.Request) {
	var req RequestSessionUsers
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	users, err := sessionManager.GetUsers(req.SessionID)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	render.JSON(w, r, ResponseUsersReq{
		Response:  resp.OK(),
		UsersReqs: users,
	})

}

func session_add_user(w http.ResponseWriter, r *http.Request) {
	var req RequestAddUser
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	err = sessionManager.AddUser(&req.User, req.SessionID)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	render.JSON(w, r, resp.OK())

}

func main() {
	model, err := model.New("http://0.0.0.0:5000/rec")
	if err != nil {
		fmt.Println("Error with model")
		os.Exit(1)
	}
	sessionManager, err = sessions.NewSessionManager("localhost:6379", "", 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// TODO : add config
	router := chi.NewRouter()
	router.Get("/cards", cards.New(model))
	router.Post("/sessions", session_create_page)
	router.Get("/sessions/{id}", session_get_data)
	router.Patch("/sessions/{id}", session_add_user)

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

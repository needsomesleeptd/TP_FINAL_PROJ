package sessions_handler

import (
	"net/http"
	auth_handler "test_backend_frontend/internal/http-server/handlers/auth"
	"test_backend_frontend/internal/lib/api/response"
	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models"
	session "test_backend_frontend/internal/sessions"
	"test_backend_frontend/pkg/auth_utils"

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

type RequestModifyUser struct {
	NewName        string    `json:"newName"`
	NewRequest     string    `json:"newRequest"`
	SessionID      uuid.UUID `json:"sessionID"`
	UserIDToModify uint64    `json:"userIDToModify"` //the id of user to modify
}

type RequestGetAllSessionsByUser struct {
	UserID uint64 `json:"userID"`
}

type ResponseGetAllSessionsByUser struct {
	Response resp.Response
	Sessions []session.Session `json:"sessions"`
}

func SessionCreatePage(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//userReq := models.NewUserReq(2, "anyname", "initializer of  a party")
		var payload *auth_utils.Payload
		cookie, err := r.Cookie(auth_handler.COOKIE_NAME)
		if err != nil {
			render.JSON(w, r, response.Error("Error with cookie"))
			return
		}
		//fmt.Print(cookie.Value)
		payload, err = sessionManager.TokenHandler.ParseToken(cookie.Value, sessionManager.Secret)
		if err != nil {
			render.JSON(w, r, response.Error("Error getting data"))
			return
		}

		var sessionName = "session Name"
		userReq := models.UserReq{ID: payload.ID, Name: payload.Login, Request: "fill me!"}
		sessionID, err := sessionManager.CreateSession(&userReq, sessionName)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, ResponseSessionID{
			Response:  resp.OK(),
			SessionID: sessionID,
		})
	}
}

func SessionGetData(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func SessionAdduser(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func SessionModifyuser(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestModifyUser

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		updateReq := models.NewUserReq(req.UserIDToModify, req.NewName, req.NewRequest)
		err = sessionManager.ModifyUser(req.SessionID, req.UserIDToModify, updateReq)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, resp.OK())

	}
}

func SessionGetUserSessions(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGetAllSessionsByUser
		var sessions []session.Session
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		sessions, err = sessionManager.GetUserSessions(req.UserID)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, ResponseGetAllSessionsByUser{
			Response: resp.OK(),
			Sessions: sessions,
		})
	}
}

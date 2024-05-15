package sessions_handler

import (
	"net/http"
	"strings"
	"test_backend_frontend/internal/lib/api/response"
	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models"
	session "test_backend_frontend/internal/sessions"
	"test_backend_frontend/pkg/auth_utils"
	"time"

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

type ResponseSession struct {
	Response resp.Response
	Session  session.Session `json:"session"`
}

type RequestSession struct {
	SessionID uuid.UUID `json:"sessionID"`
}

type RequestCreateSession struct {
	SessionName      string    `json:"sessionName"`
	SessionPeopleCap int       `json:"sessionPeopleCap"`
	SessionOwner     string    `json:"sessonOwner"`
	Description      string    `json:"description"`
	TimeEnds         time.Time `json:"timeEnds"`
}

type RequestAddUser struct {
	//Jwt       string    `json:"jwt"`
	SessionID uuid.UUID `json:"sessionID"`
}

type RequestModifyUser struct {
	NewName        string    `json:"newName"`
	NewRequest     string    `json:"newRequest"`
	SessionID      uuid.UUID `json:"sessionID"`
	UserIDToModify uint64    `json:"userIDToModify"` //the id of user to modify
	Categories     []string  `json:"newCategories"`
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
		var req RequestCreateSession
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		var payload *auth_utils.Payload
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		payload, err = sessionManager.TokenHandler.ParseToken(token, sessionManager.Secret)
		if err != nil {
			render.JSON(w, r, response.Error("Error getting data"))
			return
		}

		userReq := models.UserReq{ID: payload.ID, Name: payload.Login, Request: ""}
		//var duration time.Duration = 1e9
		sessionID, err := sessionManager.CreateSession(&userReq, req.SessionName, req.SessionPeopleCap, req.TimeEnds, req.Description)
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

func SessionsGetSessionData(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSession
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		session, err := sessionManager.GetSession(req.SessionID)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, ResponseSession{
			Response: resp.OK(),
			Session:  *session,
		})

	}
}

func SessionGetData(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSession
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
		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		var payload *auth_utils.Payload
		payload, err = sessionManager.TokenHandler.ParseToken(token, sessionManager.Secret)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
		}
		user := models.UserReq{ID: payload.ID, Name: payload.Login}
		err = sessionManager.AddUser(&user, req.SessionID)
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
		updateReq.Categories = req.Categories
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

func SessionContinueScrolling(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSession
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		session, err := sessionManager.GetSession(req.SessionID)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		if session.Status == models.Ended {
			err = sessionManager.ChangeSessionStatus(req.SessionID, models.Scrolling)
			if err != nil {
				render.JSON(w, r, response.Error(err.Error()))
				return
			}
		}
		render.JSON(w, r, resp.OK())
	}
}

func SessionDeleteUser(sessionManager *session.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSession
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var payload *auth_utils.Payload
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		payload, err = sessionManager.TokenHandler.ParseToken(token, sessionManager.Secret)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
		}

		err = sessionManager.DeletePersonFromSession(req.SessionID, payload.ID)
		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		render.JSON(w, r, resp.OK())
	}
}

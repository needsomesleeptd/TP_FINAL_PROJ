package auth_handler

import (
	"net/http"
	"test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_dto"
	auth_service "test_backend_frontend/internal/services/auth"
	"time"

	"github.com/go-chi/render"
)

const (
	COOKIE_NAME = "auth_jwt"
)

type RequestSignUp struct {
	User models_dto.User `json:"user"`
}
type RequestSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SignUp(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		candidate := models_dto.FromDtoUser(&req.User)
		err = authService.SignUp(&candidate)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func SignIn(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignIn
		var tokenStr string
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		candidate := models.User{Login: req.Login, Password: req.Password}
		tokenStr, err = authService.SignIn(&candidate)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		cookie := http.Cookie{
			Name:     COOKIE_NAME,
			Value:    tokenStr,
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: true, // Ilya, change this if you want to access cookie
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		render.JSON(w, r, response.OK())
	}
}

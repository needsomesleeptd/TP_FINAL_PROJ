package cards

import (
	"errors"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"strings"
	"test_backend_frontend/internal/models/models_dto"
	auth_service "test_backend_frontend/internal/services/auth"
	"test_backend_frontend/pkg/auth_utils"

	resp "test_backend_frontend/internal/lib/api/response"
)

type Request struct {
	Prompt    string `json:"prompt"`
	SessionID string `json:"sessionID"`
}

type Response struct {
	resp.Response
	Cards []*models_dto.Card `json:"cards"`
}

type TokenParser interface {
	ParseToken(tokenString string, key string) (*auth_utils.Payload, error)
}

type CardsSearcher interface {
	CardsSearch(prompt string, sessionId string, userId uint64) ([]*models_dto.Card, error)
}

func New(searcher CardsSearcher, tokenizer TokenParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, resp.Error("failed to get auth token"))
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		payload, err := tokenizer.ParseToken(token, auth_service.SECRET)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to parse token"))
			return
		}

		cards, err := searcher.CardsSearch(req.Prompt, req.SessionID, payload.ID)

		if err != nil {
			render.JSON(w, r, resp.Error("failed to get cards"))
			return
		}

		responseOK(w, r, cards)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, cards []*models_dto.Card) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Cards:    cards,
	})
}

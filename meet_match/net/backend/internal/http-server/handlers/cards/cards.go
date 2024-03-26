package cards

import (
	"errors"
	"github.com/go-chi/render"
	"io"
	"net/http"

	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/model"
)

type Request struct {
	Prompt       string `json:"prompt"`
	Page         int    `json:"page"`
	CardsPerPage int    `json:"cardsPerPage"`
}

type Response struct {
	resp.Response
	Cards []model.Card `json:"cards"`
}

type CardsSearcher interface {
	CardsSearch(prompt string, fromLine int, toLine int) ([]model.Card, error)
}

func New(searcher CardsSearcher) http.HandlerFunc {
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

		if req.Page <= 0 {
			render.JSON(w, r, resp.Error("wrong page"))
			return
		}
		if req.CardsPerPage <= 0 {
			render.JSON(w, r, resp.Error("wrong cards count"))
			return
		}

		fromLine := (req.Page - 1) * req.CardsPerPage
		toLine := req.Page * req.CardsPerPage
		cards, err := searcher.CardsSearch(req.Prompt, fromLine, toLine)

		if err != nil {
			render.JSON(w, r, resp.Error("failed to get cards"))
			return
		}

		responseOK(w, r, cards)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, cards []model.Card) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Cards:    cards,
	})
}

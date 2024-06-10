package scroll

import (
	"io"
	"log"
	"net/http"
	"test_backend_frontend/internal/lib/api/response"
	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/middleware/auth_middleware"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_dto"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CheckMatchRequest struct {
	SessionID string `json:"sessionID"`
}

type CardsMatchChecker interface {
	GetMatchCards(session_id uuid.UUID, userID uint64) ([]*models.Card, error)
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
}

type Response struct {
	resp.Response
	IsMatched bool               `json:"is_matched"`
	Cards     []*models_dto.Card `json:"cards"`
}

func NewCheckHandler(checker CardsMatchChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error("unable to fetch userID to check the data"))
			return
		}
		var req CheckMatchRequest
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		uid, err := uuid.Parse(req.SessionID)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to parse uuid"))
			return
		}
		cards, err := checker.GetMatchCards(uid, userID)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to get match"))
			return
		}
		var dtoCards []*models_dto.Card
		for _, v := range cards {
			dtoCards = append(dtoCards, models_dto.ToDTOCard(v))
		}
		log.Printf("sending cards %v", dtoCards)
		responseOK(w, r, dtoCards)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, cards []*models_dto.Card) {
	render.JSON(w, r, Response{
		Response:  resp.OK(),
		IsMatched: len(cards) > 0,
		Cards:     cards,
	})
}

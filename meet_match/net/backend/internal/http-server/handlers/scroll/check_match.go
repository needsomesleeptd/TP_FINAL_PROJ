package scroll

import (
	"io"
	"log"
	"net/http"
	resp "test_backend_frontend/internal/lib/api/response"
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
	GetMatchCards(session_id uuid.UUID) ([]*models.Card, error)
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
}

// TODO: one card
// TODO: jwt in header
type Response struct {
	resp.Response
	IsMatched bool               `json:"is_matched"`
	Cards     []*models_dto.Card `json:"cards"`
}

func NewCheckHandler(checker CardsMatchChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		//wasMatched,err := checker.IsMatchHappened()
		cards, err := checker.GetMatchCards(uid)
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

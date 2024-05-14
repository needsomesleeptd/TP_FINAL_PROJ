package match_handler

import (
	"errors"
	"io"
	"net/http"
	"test_backend_frontend/internal/http-server/handlers/cards"
	sessions_handler "test_backend_frontend/internal/http-server/handlers/session"
	"test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models/models_dto"
	match_service "test_backend_frontend/internal/services/match"

	"github.com/go-chi/render"
)

func GetMatchedCards(matchServ match_service.IMatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sessions_handler.RequestSession

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		cardsMatched, err := matchServ.GetMatchedCardsBySession(req.SessionID)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}

		cardsMatchedDTO := make([]*models_dto.Card, len(cardsMatched))

		for i, card := range cardsMatched {
			cardsMatchedDTO[i] = models_dto.ToDTOCard(card)
		}

		render.JSON(w, r, cards.Response{Response: response.OK(), Cards: cardsMatchedDTO})
	}
}

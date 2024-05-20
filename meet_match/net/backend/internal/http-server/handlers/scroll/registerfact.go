package scroll

import (
	"io"
	"log"
	"net/http"
	"strings"
	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_dto"
	auth_service "test_backend_frontend/internal/services/auth"
	"test_backend_frontend/pkg/auth_utils"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ScrollFactRegistrateRequest struct {
	SessionId string `json:"sessionID"`
	PlaceId   uint64 `json:"placeID"`
	IsLiked   bool   `json:"is_liked"`
}

type ScrollFactRegistrator interface {
	RegisterFact(scrolled *models.FactScrolled) error
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
}

type TokenParser interface {
	ParseToken(tokenString string, key string) (*auth_utils.Payload, error)
}

type PlaceProvider interface {
	GetCard(id uint64) (*models.Card, error)
}

func NewScrollFactRegistrateHandler(registrator ScrollFactRegistrator, tokenizer TokenParser, cardProv PlaceProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ScrollFactRegistrateRequest
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		sessionId, err := uuid.Parse(req.SessionId)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to parse uuid"))
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, resp.Error("failed to get auth token"))
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		userId, err := tokenizer.ParseToken(token, auth_service.SECRET)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to parse token"))
			return
		}

		fact := &models.FactScrolled{
			SessionId: sessionId,
			UserId:    userId.ID,
			PlacesId:  req.PlaceId,
			IsLiked:   req.IsLiked,
		}

		err = registrator.RegisterFact(fact)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to save fact"))
			log.Println(err.Error())
			return
		}
		is_match := false

		/*if req.IsLiked {
			is_match, err = registrator.IsMatchHappened(fact)
			if err != nil {
				render.JSON(w, r, resp.Error("match check issue"))
				return
			}
		}
		*/
		var places []*models_dto.Card
		/*if is_match {
			place, err := cardProv.GetCard(fact.PlacesId)
			if err != nil {
				render.JSON(w, r, resp.Error("place get issue"))
				return
			}

			places = append(places, models_dto.ToDTOCard(place))
		}*/

		render.JSON(w, r, Response{
			Response:  resp.OK(),
			IsMatched: is_match,
			Cards:     places,
		})
	}
}

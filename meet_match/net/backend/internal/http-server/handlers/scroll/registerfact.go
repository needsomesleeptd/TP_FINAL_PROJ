package scroll

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net/http"
	auth_handler "test_backend_frontend/internal/http-server/handlers/auth"
	resp "test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_dto"
	auth_service "test_backend_frontend/internal/services/auth"
	"test_backend_frontend/pkg/auth_utils"
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

func NewScrollFactRegistrateHandler(registrator ScrollFactRegistrator, tokenizer TokenParser) http.HandlerFunc {
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

		cookie, err := r.Cookie(auth_handler.COOKIE_NAME)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to get cookie"))
			return
		}

		userId, err := tokenizer.ParseToken(cookie.Value, auth_service.SECRET)
		if err != nil {
			render.JSON(w, r, resp.Error("failed to parse cookie"))
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
			return
		}

		is_match, err := registrator.IsMatchHappened(fact)
		if err != nil {
			render.JSON(w, r, resp.Error("match check issue"))
			return
		}

		var places []*models_dto.Card
		if is_match {
			// TODO: fill the gaps
			places = append(places, &models_dto.Card{
				Id:       fact.PlacesId,
				ImgUrl:   "",
				CardName: "",
				Rating:   0,
			})
		}

		render.JSON(w, r, Response{
			Response:  resp.OK(),
			IsMatched: is_match,
			Cards:     places,
		})
	}
}
package scroll_stats_handler

import (
	"net/http"
	"test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/middleware/auth_middleware"
	"test_backend_frontend/internal/models/models_dto"
	scroll_stats_serv "test_backend_frontend/internal/services/scrollStats"

	"github.com/go-chi/render"
)

type ResponseScrollStats struct {
	response.Response
	PersonStats models_dto.PersonScrollStats `json:"peron_stats"`
}

func GetUserStats(matchServ scroll_stats_serv.IScrolledStatsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error("unable to fetch userID to get feedback"))
			return
		}
		personStats, err := matchServ.GetPersonStats(userID)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		personStatsDto := models_dto.ToDToPersonScrollStats(*personStats)
		resp := ResponseScrollStats{Response: response.OK(), PersonStats: *personStatsDto}
		render.JSON(w, r, resp)

	}
}

package feedback_handler

import (
	"net/http"
	"test_backend_frontend/internal/lib/api/response"
	"test_backend_frontend/internal/middleware/auth_middleware"
	"test_backend_frontend/internal/models"
	feedback_service "test_backend_frontend/internal/services/feedback"
	"time"

	"github.com/go-chi/render"
)

type RequestSaveFeedback struct {
	Description string    `json:"description,omitempty"`
	HasGone     bool      `json:"has_gone"`
	Datetime    time.Time `json:"datetime,omitempty"`
}

func SaveFeedback(feedbackService feedback_service.IFeedbackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSaveFeedback
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error("unable to fetch userID to get feedback"))
			return
		}

		feedback := models.Feedback{
			Description: req.Description,
			HasGone:     req.HasGone,
			Datetime:    req.Datetime,
			UserID:      userID,
		}
		err = feedbackService.SaveFeedback(feedback)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
		}
		render.JSON(w, r, response.OK())
	}
}

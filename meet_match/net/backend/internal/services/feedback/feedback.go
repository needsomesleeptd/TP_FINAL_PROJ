package feedback_service

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/feedback/feedback_repo"

	"github.com/pkg/errors"
)

type IFeedbackService interface {
	SaveFeedback(feedback models.Feedback) error
}

func NewFeedbackService(repo feedback_repo.FeedbackRepository) IFeedbackService {
	return &FeedbackService{repo: repo}
}

type FeedbackService struct {
	repo feedback_repo.FeedbackRepository
}

func (s *FeedbackService) SaveFeedback(feedback models.Feedback) error {
	err := s.repo.SaveFeedBack(feedback)
	if err != nil {
		return errors.Wrap(err, "error in feedback service")
	}
	return nil
}

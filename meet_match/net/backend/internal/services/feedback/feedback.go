package feedback_service

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/feedback/feedback_repo"
	"test_backend_frontend/internal/services/scroll"
	session "test_backend_frontend/internal/sessions"

	"github.com/pkg/errors"
)

type IFeedbackService interface {
	SaveFeedback(feedback models.Feedback) error
}

func NewFeedbackService(repo feedback_repo.FeedbackRepository) IFeedbackService {
	return &FeedbackService{repo: repo}
}

type FeedbackService struct {
	repo        feedback_repo.FeedbackRepository
	scrollServ  scroll.ScrollUseCase
	sessionServ session.SessionManager
}

func (s *FeedbackService) SaveFeedback(feedback models.Feedback) error {
	err := s.repo.SaveFeedBack(feedback)
	if err != nil {
		return errors.Wrap(err, "error in feedback service")
	}
	return nil
}

func (s *FeedbackService) RequestFeedBackData(userID uint64) (*models.Card, error) {
	sessions, err := s.sessionServ.GetUserSessions(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting sessions for feedback")
	}
	idLastSession := sessions[len(sessions)-1].SessionID
	matchedCards, err := s.scrollServ.GetMatchCards(idLastSession)
	if err != nil {
		return nil, errors.Wrap(err, "error getting marchedCards for feedback")
	}
	lastMatch := matchedCards[0]
	return lastMatch, nil
}

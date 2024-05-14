package match_service

import (
	match_repo "test_backend_frontend/internal/services/match/matchRepo"
	session "test_backend_frontend/internal/sessions"
)

type IMatchService interface {
	IsFeedBackRequired(userID uint64) (bool, error)
}

type MatchService struct {
	matchRepo  match_repo.IMatchRepo
	sessionMan session.SessionManager
}

func (m *MatchService) IsFeedBackRequired(userID uint64) (bool, error) {
	sessions, err := m.sessionMan.GetUserSessions(userID)
	if err != nil {
		return false,err
	}
	
	m.matchRepo.GetMatchesNoFeedback()

}

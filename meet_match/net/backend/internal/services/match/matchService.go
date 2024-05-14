package match_service

import (
	"fmt"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/cards/repository"
	match_repo "test_backend_frontend/internal/services/match/matchRepo"
	session "test_backend_frontend/internal/sessions"

	"github.com/google/uuid"
)

type IMatchService interface {
	GetMatchedCardsBySession(sessionID uuid.UUID) ([]*models.Card, error)
}

func NewMatchService(matchRepoSrc match_repo.IMatchRepo, sessionManSc session.SessionManager, cardRepoSrc repository.CardRepository) IMatchService {
	return &MatchService{
		matchRepo:  matchRepoSrc,
		sessionMan: sessionManSc,
		cardRepo:   cardRepoSrc,
	}

}

type MatchService struct {
	matchRepo  match_repo.IMatchRepo
	sessionMan session.SessionManager
	cardRepo   repository.CardRepository
}

func (m *MatchService) GetMatchedCardsBySession(sessionID uuid.UUID) ([]*models.Card, error) {
	matches, err := m.matchRepo.GetMatchesBySession(sessionID)
	if err != nil {
		return nil, err
	}
	matchedCards := make([]*models.Card, len(matches))
	for i, match := range matches {
		matchedCards[i], err = m.cardRepo.GetCard(match.CardMatchedID)
		if err != nil {
			fmt.Printf("error extracting matchedCards by session: %w\n", err) //TODO:: add logging
		}

	}
	return matchedCards, nil
}

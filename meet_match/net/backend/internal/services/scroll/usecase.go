package scroll

import (
	"github.com/google/uuid"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/scroll/scroll_repo"
)

type ScrollUseCase interface {
	RegisterFact(scrolled *models.FactScrolled) error
	IsMatchHappened(session_id uuid.UUID) (bool, error)
}

type usecase struct {
	repo scroll_repo.ScrollRepository
}

func (u usecase) RegisterFact(scrolled *models.FactScrolled) error {
	//TODO implement me
	panic("implement me")
}

func (u usecase) IsMatchHappened(session_id uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

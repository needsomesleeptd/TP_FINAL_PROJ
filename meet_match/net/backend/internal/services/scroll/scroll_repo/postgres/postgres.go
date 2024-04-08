package postgres

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"test_backend_frontend/internal/models"
)

type scrollRepository struct {
	db *gorm.DB
}

func (s scrollRepository) AddScrollFact(fact *models.FactScrolled) error {
	//TODO implement me
	panic("implement me")
}

func (s scrollRepository) GetAllLikedPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (s scrollRepository) GetAllUsersIdsForSession(session_id uuid.UUID) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (s scrollRepository) GetAllPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

package postgres

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
)

type scrollRepository struct {
	db *gorm.DB
}

const MAX_LIMIT = 100

func (s scrollRepository) AddScrollFact(fact *models.FactScrolled) error {
	pgFact := models_da.ToPostgresFactScrolled(fact)
	tx := s.db.Create(&pgFact)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "scroll.repository.AddScrollFact error")
	}
	return nil
}

func (s scrollRepository) GetAllLikedPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error) {
	var ids []uint64

	tx := s.db.Limit(MAX_LIMIT).
		Model(&models_da.FactScrolled{}).
		Select("places_id").
		Where("session_id = ? AND user_id = ? AND is_liked = true", session_id.String(), user_id).
		Find(&ids)

	if tx.Error != nil {
		return nil,
			errors.Wrap(tx.Error, "scroll.repository.GetAllLikedPlaces error")
	}

	return ids, nil
}

func (s scrollRepository) GetAllUsersIdsForSession(session_id uuid.UUID) ([]uint64, error) {
	var ids []uint64

	tx := s.db.Limit(MAX_LIMIT).
		Model(&models_da.FactScrolled{}).
		Select("user_id").
		Where("session_id = ?", session_id.String()).
		Find(&ids)

	if tx.Error != nil {
		return nil,
			errors.Wrap(tx.Error, "scroll.repository.GetAllUsersIdsForSession error")
	}

	return ids, nil
}

func (s scrollRepository) GetAllPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

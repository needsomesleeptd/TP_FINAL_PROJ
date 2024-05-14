package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
	"test_backend_frontend/internal/services/cards/repository"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) repository.CardRepository {
	return &cardRepository{db: db}
}

func (c cardRepository) GetCard(id uint64) (*models.Card, error) {
	var pgCard *models_da.Card
	tx := c.db.First(&pgCard, "place_id = ?", id)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "card.repository.GetCard error")
	}

	if pgCard.Subway == nil && pgCard.Place != nil {
		var place *models_da.Card
		tx := c.db.Limit(1).Find(&place, "place_id = ?", pgCard.Place)
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "card.repository.GetCard error")
		}

		if place.Subway != nil {
			pgCard.Subway = place.Subway
		}
	}

	return models_da.ToModelCard(pgCard), nil
}

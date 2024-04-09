package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
)

type cardRepository struct {
	db *gorm.DB
}

func (c cardRepository) GetCard(id uint64) (*models.Card, error) {
	var pgCard *models_da.Card
	tx := c.db.First(&pgCard, "id = ?", id)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "card.repository.GetCard error")
	}

	return models_da.ToModelCard(pgCard), nil
}

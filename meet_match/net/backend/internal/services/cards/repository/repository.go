package repository

import "test_backend_frontend/internal/models"

type CardRepository interface {
	GetCard(id uint64) (*models.Card, error)
}

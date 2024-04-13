package models_da

import "test_backend_frontend/internal/models"

type Card struct {
	ID          uint64 `gorm:"column:place_id"`
	Url         string `gorm:"column:url"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

func (Card) TableName() string {
	return "places"
}

func ToModelCard(card *Card) *models.Card {
	return &models.Card{
		Id:       card.ID,
		ImgUrl:   card.Url,
		CardName: card.Title,
		Rating:   0,
	}
}

func ToPostgresCard(card *models.Card) *Card {
	return &Card{
		ID:    card.Id,
		Url:   card.ImgUrl,
		Title: card.CardName,
	}
}

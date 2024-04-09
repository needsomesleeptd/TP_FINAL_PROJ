package models_da

import "test_backend_frontend/internal/models"

type Card struct {
	ID    uint64 `grom:"column:id"`
	Url   string `grom:"column:url"`
	Title string `grom:"column:title"`
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

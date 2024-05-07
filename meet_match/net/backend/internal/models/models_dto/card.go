package models_dto

import "test_backend_frontend/internal/models"

type Card struct {
	Id          uint64 `json:"id"`
	ImgUrl      string `json:"image"`
	CardName    string `json:"title,card_name"`
	Rating      int    `json:"rating,omitempty"`
	Description string `json:"description,omitempty"`
}

func ToDTOCard(card *models.Card) *Card {
	return &Card{
		Id:          card.Id,
		ImgUrl:      card.ImgUrl,
		CardName:    card.CardName,
		Rating:      card.Rating,
		Description: card.Description,
	}
}

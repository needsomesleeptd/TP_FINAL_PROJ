package models_dto

import "test_backend_frontend/internal/models"

type Card struct {
	Id             uint64  `json:"id"`
	ImgUrl         string  `json:"image"`
	CardName       string  `json:"title,card_name"`
	Rating         *uint64 `json:"rating,omitempty"`
	Description    *string `json:"description,omitempty"`
	Subway         *string `json:"subway,omitempty"`
	Cost           *string `json:"cost,omitempty"`
	Timetable      *string `json:"timetable,omitempty"`
	AgeRestriction *string `json:"age_restriction,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	SiteUrl        *string `json:"site_url,omitempty"`
}

func ToDTOCard(card *models.Card) *Card {
	return &Card{
		Id:             card.Id,
		ImgUrl:         card.ImgUrl,
		CardName:       card.CardName,
		Rating:         card.Rating,
		Description:    card.Description,
		Subway:         card.Subway,
		Cost:           card.Cost,
		Timetable:      card.Timetable,
		AgeRestriction: card.AgeRestriction,
		Phone:          card.Phone,
		SiteUrl:        card.SiteUrl,
	}
}

package models_da

import "test_backend_frontend/internal/models"

type Card struct {
	ID             uint64  `gorm:"column:place_id"`
	Url            string  `gorm:"column:url"`
	Title          string  `gorm:"column:title"`
	Description    *string `gorm:"column:description"`
	FavoritesCount *uint64 `gorm:"column:favorites_count"`
	Subway         *string `gorm:"column:subway"`
	Cost           *string `gorm:"column:price"`
	Timetable      *string `gorm:"column:timetable"`
	AgeRestriction *string `gorm:"column:age_restriction"`
	Phone          *string `gorm:"column:phone"`
	SiteUrl        *string `gorm:"column:foreign_url"`

	Place *uint64 `gorm:"column:place"`
	Dates *string `gorm:"column:dates"`
}

func (Card) TableName() string {
	return "places"
}

func ToModelCard(card *Card) *models.Card {
	return &models.Card{
		Id:             card.ID,
		ImgUrl:         card.Url,
		CardName:       card.Title,
		Rating:         card.FavoritesCount,
		Description:    card.Description,
		Subway:         card.Subway,
		Cost:           card.Cost,
		Timetable:      card.Timetable,
		AgeRestriction: card.AgeRestriction,
		Phone:          card.AgeRestriction,
		SiteUrl:        card.SiteUrl,
	}
}

func ToPostgresCard(card *models.Card) *Card {
	return &Card{
		ID:             card.Id,
		Url:            card.ImgUrl,
		Title:          card.CardName,
		Description:    card.Description,
		FavoritesCount: card.Rating,
		Subway:         card.Subway,
		Cost:           card.Cost,
		Timetable:      card.Timetable,
		AgeRestriction: card.AgeRestriction,
		Phone:          card.Phone,
		SiteUrl:        card.SiteUrl,
	}
}

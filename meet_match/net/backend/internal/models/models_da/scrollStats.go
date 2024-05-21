package models_da

import "test_backend_frontend/internal/models"

type CardStats struct {
	CardID      uint64 `gorm:"column:card_id"`
	SwipedTimes uint64 `gorm:"column:swiped_times"`
}

type PesonScrollStats struct {
	PersonalStats     PersonalScrollStats `gorm:"embedded"`
	SessionsCount     uint64              `gorm:"column:sessions_count"`
	MostDislikedPlace CardStats           `gorm:"embedded;embeddedPrefix:most_disliked_place_"`
	MostLikedPlace    CardStats           `gorm:"embedded;embeddedPrefix:most_liked_place_"`
}

type PersonalScrollStats struct {
	Swiped         uint64
	PoisitveSwipes uint64
	NegativeSwipes uint64
}

func CardsStatsFromDa(stats CardStats) *models.CardStats {
	return &models.CardStats{
		CardID:      stats.CardID,
		SwipedTimes: stats.SwipedTimes,
	}
}

package models_dto // stands for data_transfer_objects

import "test_backend_frontend/internal/models"

type CardStats struct {
	CardID      uint64 `json:"card_id"`
	SwipedTimes uint64 `json:"swiped_times"`
}

type PersonScrollStats struct {
	PersonalStats        PersonalScrollStats `json:"personal_stats"`
	SessionsCount        uint64              `json:"sessions_count"`
	MostDislikedPlace    Card                `json:"most_disliked_place"`
	MostLikedPlace       Card                `json:"most_liked_place"`
	MostLikedScrolled    uint64              `json:"most_liked_scrolled_count"`
	MostDislikedScrolled uint64              `json:"most_disliked_scrolled_count"`
}

func ToDToPersonScrollStats(stats models.PersonScrollStats) *PersonScrollStats {
	return &PersonScrollStats{
		PersonalStats:        PersonalScrollStats(stats.PersonalStats),
		MostLikedPlace:       Card(stats.MostLikedPlace),
		MostLikedScrolled:    stats.MostlikedScrolled,
		MostDislikedScrolled: stats.MostDislikedScrolled,
		MostDislikedPlace:    Card(stats.MostDislikedPlace),
		SessionsCount:        stats.SessionsCount,
	}

}

func FromDToPersonScrollStats(stats PersonScrollStats) *models.PersonScrollStats {
	return &models.PersonScrollStats{
		PersonalStats:        models.PersonalScrollStats(stats.PersonalStats),
		MostLikedPlace:       models.Card(stats.MostLikedPlace),
		MostlikedScrolled:    stats.MostLikedScrolled,
		MostDislikedScrolled: stats.MostDislikedScrolled,
		MostDislikedPlace:    models.Card(stats.MostDislikedPlace),
	}

}

type PersonalScrollStats struct {
	Swiped         uint64 `json:"swiped"`
	PoisitveSwipes uint64 `json:"positive_swipes"`
	NegativeSwipes uint64 `json:"negative_swipes"`
}

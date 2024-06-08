package scroll_stats_repo

import (
	"errors"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"

	"gorm.io/gorm"
)

type IScrollStatsRepo interface {
	GetMostLikedCardStats(userID uint64) (*models.CardStats, error)
	GetMostDislikedCardStats(userID uint64) (*models.CardStats, error)
	GetPersonalScrolledStats(userID uint64) (*models.PersonalScrollStats, error)
}

type ScrollStatsRepository struct {
	db *gorm.DB
}

func NewScrollRepository(db *gorm.DB) IScrollStatsRepo {
	return &ScrollStatsRepository{db: db}
}

func (s *ScrollStatsRepository) GetMostLikedCardStats(userID uint64) (*models.CardStats, error) {
	var cardStats models_da.CardStats
	err := s.db.Model(&models_da.FactScrolled{}).
		Select("place_id as card_id, COUNT(*) as swiped_times").
		Where("user_id = ?", userID).
		Where("is_liked = ?", true).
		Group("place_id").
		Order("swiped_times DESC").
		Limit(1).
		Scan(&cardStats).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound { //TODO::Fix this
		cardStats.CardID = 0
	}
	return models_da.CardsStatsFromDa(cardStats), nil
}
func (s *ScrollStatsRepository) GetMostDislikedCardStats(userID uint64) (*models.CardStats, error) {
	var cardStats models_da.CardStats
	err := s.db.Model(&models_da.FactScrolled{}).
		Select("place_id as card_id, COUNT(*) as swiped_times").
		Where("user_id = ?", userID).
		Where("is_liked = ?", false).
		Group("place_id").
		Order("swiped_times DESC").
		Limit(1).
		Scan(&cardStats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound { //TODO::Fix this
		cardStats.CardID = 0
	}
	return models_da.CardsStatsFromDa(cardStats), nil
}

func (s *ScrollStatsRepository) GetPersonalScrolledStats(userID uint64) (*models.PersonalScrollStats, error) {

	var likedCount int64
	var dislikedCount int64

	err := s.db.Model(&models_da.FactScrolled{}).Where("user_id = ? AND is_liked = ?", userID, true).Count(&likedCount).Error
	if err != nil {
		return nil, err
	}
	err = s.db.Model(&models_da.FactScrolled{}).Where("user_id = ? AND is_liked = ?", userID, false).Count(&dislikedCount).Error
	if err != nil {
		return nil, err
	}
	overallCount := likedCount + dislikedCount
	personallStats := models.PersonalScrollStats{
		Swiped:         uint64(overallCount),
		PoisitveSwipes: uint64(likedCount),
		NegativeSwipes: uint64(dislikedCount),
	}
	return &personallStats, nil
}

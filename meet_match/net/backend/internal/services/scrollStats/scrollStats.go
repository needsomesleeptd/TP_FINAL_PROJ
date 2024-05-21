package scroll_stats_serv

import (
	"log"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/cards/repository"
	scroll_stats_repo "test_backend_frontend/internal/services/scrollStats/scrollStatsRepo"
	session "test_backend_frontend/internal/sessions"

	"github.com/pkg/errors"
)

type IScrolledStatsService interface {
	GetPersonStats(userID uint64) (*models.PersonScrollStats, error)
}

type ScrolledStatsService struct {
	repoStats  scroll_stats_repo.IScrollStatsRepo
	repoCards  repository.CardRepository
	sessionMan *session.SessionManager
}

func NewScrolledStatsService(repoStatsSrc scroll_stats_repo.IScrollStatsRepo, repoCardsSrc repository.CardRepository, sessionSrc *session.SessionManager) IScrolledStatsService {
	return &ScrolledStatsService{
		repoStats:  repoStatsSrc,
		repoCards:  repoCardsSrc,
		sessionMan: sessionSrc,
	}
}

func (serv *ScrolledStatsService) GetPersonStats(userID uint64) (*models.PersonScrollStats, error) {

	var mostLikedCardStats *models.CardStats
	var mostDislikedCardStats *models.CardStats

	var mostDislikedCard *models.Card
	var mostLikedCard *models.Card

	var sessions []session.Session

	personalStats, err := serv.repoStats.GetPersonalScrolledStats(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting personal stats")
	}

	mostLikedCardStats, err = serv.repoStats.GetMostLikedCardStats(userID)

	if err != nil {
		return nil, errors.Wrap(err, "error in getting mostLiked card Stats")
	}

	mostDislikedCardStats, err = serv.repoStats.GetMostDislikedCardStats(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting mostLiked card Stats")
	}
	log.Printf("mostlikedCardStats %v", *mostLikedCardStats)
	log.Printf("mostDislikedCardStats %v", *mostDislikedCardStats)

	mostLikedCard, err = serv.repoCards.GetCard(mostLikedCardStats.CardID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting mostLiked card")
	}

	mostDislikedCard, err = serv.repoCards.GetCard(mostDislikedCardStats.CardID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting mostLiked card")
	}
	sessions, err = serv.sessionMan.GetUserSessions(userID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting user sessions")
	}
	sessionsCnt := len(sessions)
	personStats := models.PersonScrollStats{
		PersonalStats:        *personalStats,
		MostDislikedPlace:    *mostDislikedCard,
		MostLikedPlace:       *mostLikedCard,
		MostlikedScrolled:    mostLikedCardStats.SwipedTimes,
		MostDislikedScrolled: mostDislikedCardStats.SwipedTimes,
		SessionsCount:        uint64(sessionsCnt),
	}
	return &personStats, nil

}

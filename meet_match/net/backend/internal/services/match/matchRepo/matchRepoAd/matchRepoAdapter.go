package match_repo_adap

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
	match_repo "test_backend_frontend/internal/services/match/matchRepo"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MatchRepoAdapter struct {
	db *gorm.DB
}

func NewFeedbackRepo(db *gorm.DB) match_repo.IMatchRepo {
	return &MatchRepoAdapter{db: db}
}

func (r *MatchRepoAdapter) SaveMatch(match models.Match) error {
	matchDa := models_da.TODaMatch(match)
	err := r.db.Create(&matchDa).Error
	if err != nil {
		return errors.Wrap(err, "error saving match")
	}
	return err
}

func (r *MatchRepoAdapter) UpdateMatch(id uint64, match models.Match) error {
	matchDa := models_da.TODaMatch(match)
	err := r.db.Where("id = ?", id).Updates(matchDa).Error
	if err != nil {
		return errors.Wrap(err, "error updating")
	}
	return err
}

func (r *MatchRepoAdapter) GetMatchesNoFeedback(sessionID uuid.UUID) ([]models.Match, error) {
	var matchDaSlice []models_da.Match

	err := r.db.Model(&models_da.Match{}).
		Where("session_id = ?", sessionID).
		Where("got_feedback = ?", false).
		Find(&matchDaSlice).Error

	if err != nil {
		return nil, errors.Wrap(err, "error getting matches without feedback")
	}
	matchSlice := models_da.FromDaMatchSlice(matchDaSlice)
	return matchSlice, nil
}

func (r *MatchRepoAdapter) GetMatchesBySession(sessionID uuid.UUID) ([]models.Match, error) {
	var matchDaSlice []models_da.Match

	err := r.db.Model(&models_da.Match{}).Where("session_id = ?", sessionID).Find(&matchDaSlice).Error
	if err != nil {
		return nil, errors.Wrap(err, "error getting matches by session")
	}
	matchSlice := models_da.FromDaMatchSlice(matchDaSlice)
	return matchSlice, nil
}

func (r *MatchRepoAdapter) MarkMatchesAsGottenFeedback(sessionID uuid.UUID) error { // there might a race here with numerous users
	var meetups []models.Match
	if result := r.db.Model(&models.Match{}).Where("session_id = ?", sessionID).Find(&meetups); result.Error != nil {
		return result.Error
	}
	for _, meetup := range meetups {
		meetup.GotFeedback = true

		if result := r.db.Save(&meetup); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

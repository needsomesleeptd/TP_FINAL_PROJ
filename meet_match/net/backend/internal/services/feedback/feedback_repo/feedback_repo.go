package feedback_repo

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"

	"gorm.io/gorm"
)

// 2 tired to create abstract repo here

func NewFeedbackRepo(db *gorm.DB) FeedbackRepository {
	return FeedbackRepository{db: db}
}

type FeedbackRepository struct {
	db *gorm.DB
}

func (r *FeedbackRepository) SaveFeedBack(feedback models.Feedback) error {
	feedbackDA := models_da.ToDAFeedback(feedback)
	err := r.db.Create(&feedbackDA).Error
	return err
}

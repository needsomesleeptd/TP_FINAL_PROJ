package scroll

import (
	"github.com/pkg/errors"
	"slices"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/scroll/scroll_repo"
)

type ScrollUseCase interface {
	RegisterFact(scrolled *models.FactScrolled) error
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
}

type usecase struct {
	repo scroll_repo.ScrollRepository
}

func NewScrollUseCase(repository scroll_repo.ScrollRepository) ScrollUseCase {
	return &usecase{repo: repository}
}

func (u *usecase) RegisterFact(scrolled *models.FactScrolled) error {
	err := u.repo.AddScrollFact(scrolled)
	if err != nil {
		return errors.Wrap(err, "scroll.RegisterFact error")
	}

	return nil
}

func (u *usecase) IsMatchHappened(scrolled *models.FactScrolled) (bool, error) {
	userIds, err := u.repo.GetAllUsersIdsForSession(scrolled.SessionId)
	if err != nil {
		return false, errors.Wrap(err, "scroll.IsMatchHappened error")
	}

	isMatched := true
	for _, v := range userIds {
		if v != scrolled.UserId {
			likedPlaces, err := u.repo.GetAllLikedPlaces(scrolled.SessionId, v)
			if err != nil {
				return false, errors.Wrap(err, "scroll.IsMatchHappened error")
			}
			if !slices.Contains(likedPlaces, scrolled.PlacesId) {
				//TODO: добавить мажоритарное голосование, для этого нужно знать сколько челов
				// в сессии и заменить isMatched на int, считать количество матчей
				isMatched = false
				break
			}
		}
	}

	return isMatched, nil
}

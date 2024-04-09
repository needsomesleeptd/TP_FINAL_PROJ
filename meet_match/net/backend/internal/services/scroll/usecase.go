package scroll

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"slices"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/scroll/scroll_repo"
)

type ScrollUseCase interface {
	RegisterFact(scrolled *models.FactScrolled) error
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
	GetMatchCards(sessionId uuid.UUID) ([]*models.Card, error)
}

type useсase struct {
	repo scroll_repo.ScrollRepository
}

func NewScrollUseCase(repository scroll_repo.ScrollRepository) ScrollUseCase {
	return &useсase{repo: repository}
}

func (u *useсase) RegisterFact(scrolled *models.FactScrolled) error {
	err := u.repo.AddScrollFact(scrolled)
	if err != nil {
		return errors.Wrap(err, "scroll.RegisterFact error")
	}

	return nil
}

func (u *useсase) IsMatchHappened(scrolled *models.FactScrolled) (bool, error) {
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

func (u *useсase) GetMatchCards(session_id uuid.UUID) ([]*models.Card, error) {
	userIds, err := u.repo.GetAllUsersIdsForSession(session_id)
	if err != nil {
		return nil, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	var matchedIds []uint64
	var likesForAll [][]uint64

	for _, v := range userIds {
		likes, err := u.repo.GetAllLikedPlaces(session_id, v)
		if err != nil {
			return nil, errors.Wrap(err, "scroll.GetMatchCards error")
		}

		likesForAll = append(likesForAll, likes)
	}

	if len(likesForAll) <= 0 {
		return nil, nil
	}

	for i := 0; i < len(likesForAll[0]); i++ {
		isMatched := true
		for j := 1; j < len(likesForAll); j++ {
			if !slices.Contains(likesForAll[j], likesForAll[0][i]) {
				//TODO: добавить мажоритарное голосование, для этого нужно знать сколько челов
				// в сессии и заменить isMatched на int, считать количество матчей
				isMatched = false
				break
			}
		}
		if isMatched {
			matchedIds = append(matchedIds, likesForAll[0][i])
		}
	}

	// TODO: fill all gaps add repo
	var retCards []*models.Card
	for _, v := range matchedIds {
		retCards = append(retCards, &models.Card{
			Id:       v,
			ImgUrl:   "",
			CardName: "",
			Rating:   0,
		})
	}

	return retCards, nil
}

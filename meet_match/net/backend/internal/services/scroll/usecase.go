package scroll

import (
	"slices"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/cards/repository"
	match_repo "test_backend_frontend/internal/services/match/matchRepo"
	"test_backend_frontend/internal/services/scroll/scroll_repo"
	session "test_backend_frontend/internal/sessions"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ScrollUseCase interface {
	RegisterFact(scrolled *models.FactScrolled) error
	IsMatchHappened(scrolled *models.FactScrolled) (bool, error)
	GetMatchCards(sessionId uuid.UUID) ([]*models.Card, error)
}

type useсase struct {
	repo        scroll_repo.ScrollRepository
	sessionServ *session.SessionManager
	cardRepo    repository.CardRepository
	matchRepo   match_repo.IMatchRepo
}

func NewScrollUseCase(repository scroll_repo.ScrollRepository, sessMgr *session.SessionManager, cardRepo repository.CardRepository, matchRepo match_repo.IMatchRepo) ScrollUseCase {
	return &useсase{repo: repository, sessionServ: sessMgr, cardRepo: cardRepo, matchRepo: matchRepo}
}

func (u *useсase) RegisterFact(scrolled *models.FactScrolled) error {
	err := u.repo.AddScrollFact(scrolled)
	if err != nil {
		return errors.Wrap(err, "scroll.RegisterFact error")
	}
	hasHappened, err := u.IsMatchHappened(scrolled)
	if err != nil {
		return errors.Wrap(err, "scroll.RegisterFact error")
	}
	match := models.Match{
		SessionID:     scrolled.SessionId,
		Datetime:      time.Now(),
		GotFeedback:   false,
		CardMatchedID: scrolled.PlacesId,
	}
	if hasHappened {
		err := u.matchRepo.SaveMatch(match)
		if err != nil {
			return errors.Wrap(err, "scroll.RegisterFact error")
		}
	}
	return nil
}

func (u *useсase) IsMatchHappened(scrolled *models.FactScrolled) (bool, error) {
	userIds, err := u.repo.GetAllUsersIdsForSession(scrolled.SessionId)
	if err != nil {
		return false, errors.Wrap(err, "scroll.IsMatchHappened error")
	}

	matches, err := u.matchRepo.GetMatchesBySession(scrolled.SessionId)
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

			alreadyMatched := slices.ContainsFunc(matches, func(match models.Match) bool {
				return match.CardMatchedID == scrolled.PlacesId
			})

			if !slices.Contains(likedPlaces, scrolled.PlacesId) || alreadyMatched {
				//TODO: добавить мажоритарное голосование, для этого нужно знать сколько челов
				// в сессии и заменить isMatched на int, считать количество матчей
				isMatched = false
				break
			}
		}
	}

	getUsers, err := u.sessionServ.GetUsers(scrolled.SessionId)
	if err != nil {
		return false, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	if len(userIds) < len(getUsers) {
		return false, nil
	}
	if isMatched {
		err := u.sessionServ.ChangeSessionStatus(scrolled.SessionId, models.Ended)
		if err != nil {
			return false, errors.Wrap(err, "scroll.GetMatchCards session status change error")
		}
	}

	return isMatched, nil
}

func (u *useсase) GetMatchCards(session_id uuid.UUID) ([]*models.Card, error) { // notive that the matched cards are sorted by 1 user in dsc order
	// We can use matchRepo now but it works so be it
	matches, err := u.matchRepo.GetMatchesNoFeedback(session_id)

	if err != nil {
		return nil, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	var retCards []*models.Card
	for _, v := range matches {

		err := u.matchRepo.UpdateMatch(v.ID, models.Match{GotFeedback: true})

		if err != nil {
			return nil, errors.Wrap(err, "scroll.GetMatchCards error")
		}

		card, err := u.cardRepo.GetCard(v.CardMatchedID)
		if err != nil {
			return nil, errors.Wrap(err, "scroll.GetMatchCards error")
		}

		retCards = append(retCards, card)
	}

	return retCards, nil
}

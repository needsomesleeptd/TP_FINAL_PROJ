package scroll

import (
	"log"
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
	GetMatchCards(session_id uuid.UUID, userID uint64) ([]*models.Card, error)
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
	log.Printf("hasHappened value is %v", hasHappened)
	match := models.Match{
		SessionID:     scrolled.SessionId,
		Datetime:      time.Now(),
		GotFeedback:   false,
		CardMatchedID: scrolled.PlacesId,
		UserID:        scrolled.UserId,
		MatchViewed:   false,
	}
	if hasHappened { //save matches for all users
		UserIDs, err := u.repo.GetAllUsersIdsForSession(scrolled.SessionId)
		if err != nil {
			return errors.Wrap(err, "scroll.RegisterFact error getting userIds by session")
		}
		for _, UserID := range UserIDs {
			match.UserID = UserID
			err := u.matchRepo.SaveMatch(match)
			if err != nil {
				return errors.Wrap(err, "scroll.RegisterFact error")
			}
		}
	}

	return nil
}

func (u *useсase) IsMatchHappened(scrolled *models.FactScrolled) (bool, error) {
	userIds, err := u.repo.GetAllUsersIdsForSession(scrolled.SessionId)
	if err != nil {
		return false, errors.Wrap(err, "scroll.GetMatchCards error")
	}
	matches, err := u.matchRepo.GetUserMatchesBySession(scrolled.SessionId, scrolled.UserId)

	if err != nil {
		return false, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	var likesForAll [][]uint64
	for _, v := range userIds {
		likes, err := u.repo.GetAllLikedPlaces(scrolled.SessionId, v)
		if err != nil {
			return false, errors.Wrap(err, "scroll.GetMatchCards error")
		}
		likesForAll = append(likesForAll, likes)

	}

	if len(likesForAll) <= 0 {
		return false, nil
	}

	getUsers, err := u.sessionServ.GetUsers(scrolled.SessionId)
	if err != nil {
		return false, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	if len(likesForAll) < len(getUsers) {
		return false, nil
	}
	log.Printf("started seatching likes\n")
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

		alreadyLiked := slices.ContainsFunc(matches, func(match models.Match) bool {
			return likesForAll[0][i] == match.CardMatchedID
		})
		log.Printf("already liked value for %v is %v", likesForAll[0][i], alreadyLiked)

		if isMatched && !alreadyLiked {
			err := u.sessionServ.ChangeSessionStatus(scrolled.SessionId, models.Ended)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}

func (u *useсase) GetMatchCards(session_id uuid.UUID, userID uint64) ([]*models.Card, error) { // notive that the matched cards are sorted by 1 user in dsc order
	// We can use matchRepo now but it works so be it
	matches, err := u.matchRepo.GetMatchesNotViewedByUser(session_id, userID)
	log.Print("matches not viewed", matches)
	if err != nil {
		return nil, errors.Wrap(err, "scroll.GetMatchCards error")
	}

	var retCards []*models.Card
	for _, v := range matches {

		err := u.matchRepo.UpdateMatch(v.ID, models.Match{MatchViewed: true})

		if err != nil {
			return nil, errors.Wrap(err, "scroll.GetMatchCards error")
		}

		card, err := u.cardRepo.GetCard(v.CardMatchedID)
		if err != nil {
			return nil, errors.Wrap(err, "scroll.GetMatchCards error")
		}

		retCards = append(retCards, card)
	}
	log.Printf("returning cards  for sending %v", retCards)
	return retCards, nil
}

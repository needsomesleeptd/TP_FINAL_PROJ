package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/pkg/auth_utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	SessionID   uuid.UUID            `json:"sessionID" redis:"SessionID"`
	SessionName string               `json:"sessionName" redis:"sessionName"`
	Users       []models.UserReq     `json:"users" redis:"users"`
	MaxPeople   int                  `json:"maxPeople"  redis:"maxPeople"`
	Status      models.SessionStatus `json:"status" redis:"status"`
	TimeEnds    time.Time            `json:"timeEnds" redis:"timeEnds"`
	Description string               `json:"description" redis:"description"`
}

type SessionManager struct {
	Client       *redis.Client
	Secret       string
	TokenHandler auth_utils.ITokenHandler
	SessionIDs   []uuid.UUID
}

func NewSessionManager(addr, password string, db int, tokenHandler auth_utils.ITokenHandler, secret string) (*SessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &SessionManager{Client: client, Secret: secret, TokenHandler: tokenHandler}, nil
}

func (s *SessionManager) CreateSession(creator *models.UserReq, sessionName string, peopleCap int, timeEnds time.Time, description string) (uuid.UUID, error) {
	newSessionID := uuid.New()

	s.SessionIDs = append(s.SessionIDs, newSessionID)
	session := Session{
		SessionID:   newSessionID,
		SessionName: sessionName,
		Users:       []models.UserReq{*creator},
		MaxPeople:   peopleCap,
		TimeEnds:    timeEnds,
		Status:      models.Waiting,
		Description: description,
	}
	session.TimeEnds = session.TimeEnds.Add(time.Hour * 200)
	marhsalledData, err := json.Marshal(session)
	if err != nil {
		return uuid.Max, errors.New("failed to marshall Session")
	}
	err = s.Client.Set(context.TODO(), newSessionID.String(), marhsalledData, 0).Err()
	if err != nil {
		return uuid.Max, err
	}

	return newSessionID, err
}

func (s *SessionManager) AddUser(user *models.UserReq, sessionID uuid.UUID) error {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return errors.Join(errors.New("add user error"), err)
	}
	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return errors.Join(errors.New("add user error"), err)
	}

	if session.Status != models.Waiting {
		return errors.New("session has already started")
	}

	isPresentInSession := slices.ContainsFunc(session.Users, func(userInSession models.UserReq) bool {
		return userInSession.ID == user.ID
	})
	if isPresentInSession {
		return errors.Join(errors.New("user is already present in session"))
	}
	session.Users = append(session.Users, *user)

	marhsalledData, err := json.Marshal(session)
	if err != nil {
		return errors.New("failed to marshall Session")
	}
	err = s.Client.Set(context.TODO(), sessionID.String(), marhsalledData, 0).Err() //TODO:: add duration here
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionManager) GetUsers(sessionID uuid.UUID) ([]models.UserReq, error) {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return nil, errors.Join(errors.New("get user error"), err)
	}
	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return nil, errors.Join(errors.New("get user error"), err)
	}

	return session.Users, nil
}

func (s *SessionManager) GetSession(sessionID uuid.UUID) (*Session, error) {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return nil, errors.Join(errors.New("get session error"), err)
	}
	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return nil, errors.Join(errors.New("get session error"), err)
	}
	if session.TimeEnds.UTC().Before(time.Now().UTC()) {
		session.Status = models.Scrolling
	}

	return &session, nil
}

func (s *SessionManager) ModifyUser(sessionID uuid.UUID, userModifyID uint64, user *models.UserReq) error {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return errors.Join(errors.New("modify user error"), err)
	}

	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return errors.Join(errors.New("modify user error"), err)
	}

	for i, userSession := range session.Users {

		if userSession.ID == userModifyID {
			session.Users[i] = *user

			someoneNotEnteredReq := slices.ContainsFunc(session.Users, func(userInSession models.UserReq) bool {
				return userInSession.Request == ""
			})
			if len(session.Users) >= session.MaxPeople && !someoneNotEnteredReq {
				session.Status = models.Scrolling
			}
			marhsalledData, err := json.Marshal(session)
			if err != nil {
				return errors.New("failed to marshall Session")
			}
			err = s.Client.Set(context.TODO(), sessionID.String(), marhsalledData, 0).Err() //add session TimeHere
			if err != nil {
				return err
			}

			return nil
		}

	}
	return fmt.Errorf("haven't found  users with this ID")
}

// Не забудтьте отчистить хранилище реддис

func (s *SessionManager) GetUserSessions(userID uint64) ([]Session, error) {
	// s.Client.FlushAll(context.TODO())
	keys, err := s.Client.Keys(context.TODO(), "*").Result()
	if err != nil {
		return nil, errors.Join(errors.New("getting keys"), err)
	}
	var session Session
	var sessions []Session
	for _, sessionID := range keys {

		sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID).Result()
		if err != nil && err != redis.Nil {
			return nil, errors.Join(errors.New("getting session error"), err)
		}
		err = json.Unmarshal([]byte(sessionMarshalled), &session)
		if err != nil {
			//return nil, errors.Join(errors.New("getting session error"), err)
			continue
		}
		for _, user := range session.Users {
			if userID == user.ID {
				sessions = append(sessions, session)
				break
			}
		}
	}
	return sessions, nil
}

func (s *SessionManager) ChangeSessionStatus(sessionID uuid.UUID, status models.SessionStatus) error {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return errors.Join(errors.New("changing user status error"), err)
	}
	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return errors.Join(errors.New("changing user status error"), err)
	}

	session.Status = status

	marhsalledData, err := json.Marshal(session)
	if err != nil {
		return errors.New("failed to marshall Session")
	}
	err = s.Client.Set(context.TODO(), sessionID.String(), marhsalledData, 0).Err() //TODO:: add duration here
	if err != nil {
		return errors.Join(errors.New("changing user status error"), err)
	}
	return nil
}

func (s *SessionManager) DeletePersonFromSession(sessionID uuid.UUID, userID uint64) error {
	var session Session
	sessionMarshalled, err := s.Client.Get(context.TODO(), sessionID.String()).Result()
	if err != nil {
		return errors.Join(errors.New("deleting user error"), err)
	}
	err = json.Unmarshal([]byte(sessionMarshalled), &session)
	if err != nil {
		return errors.Join(errors.New("deleting user error"), err)
	}
	index := slices.IndexFunc(session.Users, func(userInSession models.UserReq) bool {
		return userInSession.ID == userID
	})

	// we haven't found a person
	if index == -1 {
		return errors.New("the person doesn't exist in this session")
	}
	//deleting a person
	session.Users = slices.Delete(session.Users, index, index+1)

	marhsalledData, err := json.Marshal(session)
	if err != nil {
		return errors.New("failed to marshall Session")
	}
	if len(session.Users) > 0 {
		err = s.Client.Set(context.TODO(), sessionID.String(), marhsalledData, 0).Err() //TODO:: add duration here
	} else {
		err = s.Client.Del(context.Background(), sessionID.String()).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

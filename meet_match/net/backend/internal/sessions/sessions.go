package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/pkg/auth_utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	SessionID    uuid.UUID        `json:"sessionID" redis:"SessionID"`
	SessionName  string           `json:"sessionName" redis:"sessionName"`
	Users        []models.UserReq `json:"users" redis:"users"`
	MaxPeople    int              `json:"maxPeople"  redis:"maxPeople"`
	HasStarted   bool             `json:"hasStarted" redis:"hasStarted"`
	TimeDuration time.Duration    `json:"duration" redis:"duration"`
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

func (s *SessionManager) CreateSession(creator *models.UserReq, sessionName string, peopleCap int, timeDur time.Duration) (uuid.UUID, error) {
	newSessionID := uuid.New()

	s.SessionIDs = append(s.SessionIDs, newSessionID)
	session := Session{
		SessionID:    newSessionID,
		SessionName:  sessionName,
		Users:        []models.UserReq{*creator},
		MaxPeople:    peopleCap,
		TimeDuration: timeDur,
		HasStarted:   false,
	}
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
	for i, user := range session.Users {

		if user.ID == userModifyID {
			session.Users[i] = user
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
			return nil, errors.Join(errors.New("getting session error"), err)
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

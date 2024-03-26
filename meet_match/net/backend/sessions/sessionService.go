package session

import (
	"context"
	"encoding/json"
	"meetMatch/models"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	Client     *redis.Client
	SessionIDs []uuid.UUID
}

func NewSessionManager(addr, password string, db int) (*SessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &SessionManager{Client: client}, nil
}

func (s *SessionManager) CreateSession(creator *models.UserReq) (uuid.UUID, *redis.StatusCmd) {
	newSessionID := uuid.New()
	s.SessionIDs = append(s.SessionIDs, newSessionID)
	marhsalledData, err := json.Marshal(*creator)
	if err != nil {
		panic("Error in marshalling")
	}
	status := s.Client.Set(context.TODO(), newSessionID.String(), marhsalledData, 3600*time.Second)
	return newSessionID, status
}

func (s *SessionManager) AddUser(user *models.UserReq, sessionID uuid.UUID) *redis.StatusCmd {
	marhsalledData, err := json.Marshal(*user)
	if err != nil {
		panic("Error in marshalling")
	}
	status := s.Client.Set(context.TODO(), sessionID.String(), marhsalledData, 3600*time.Second) //TODO:: fix add expiration + error
	return status
}

func (s *SessionManager) GetUsers(sessionID uuid.UUID) ([]models.UserReq, error) {
	var users []models.UserReq
	err := s.Client.HGetAll(context.TODO(), sessionID.String()).Scan(&users)
	return users, err
}

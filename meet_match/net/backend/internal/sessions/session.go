package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"test_backend_frontend/internal/models"

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

func (s *SessionManager) CreateSession(creator *models.UserReq) (uuid.UUID, error) {
	newSessionID := uuid.New()
	s.SessionIDs = append(s.SessionIDs, newSessionID)
	marhsalledData, err := json.Marshal(*creator)
	if err != nil {
		return uuid.Max, errors.New("failed to marshall Session")
	}
	st, err := s.Client.RPush(context.TODO(), newSessionID.String(), marhsalledData).Result()
	fmt.Println(st)
	if err != nil {
		return uuid.Max, err
	}
	return newSessionID, err
}

func (s *SessionManager) AddUser(user *models.UserReq, sessionID uuid.UUID) error {
	marhsalledData, err := json.Marshal(*user)
	if err != nil {
		return err
	}
	st, err := s.Client.LPush(context.TODO(), sessionID.String(), marhsalledData).Result() //TODO:: fix add expiration + error
	fmt.Print(st)
	return err
}

func (s *SessionManager) GetUsers(sessionID uuid.UUID) ([]models.UserReq, error) {
	var users []models.UserReq
	list, err := s.Client.LRange(context.TODO(), sessionID.String(), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		var user models.UserReq
		err := json.Unmarshal([]byte(item), &user) // Unmarshal Redis list item into struct
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}

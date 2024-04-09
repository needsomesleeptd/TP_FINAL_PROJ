package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/services/auth/user_repo"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	SessionID   uuid.UUID        `json:"sessionID"`
	SessionName string           `json:"sessionName"`
	Users       []models.UserReq `json:"users"`
}

type SessionManager struct {
	Client     *redis.Client
	UserRepo   user_repo.IUserRepository
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

func (s *SessionManager) CreateSession(creator *models.UserReq, sessionName string) (uuid.UUID, error) {
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
	//var expireDuration time.Duration
	//expireDuration = 1e9

	//err = s.Client.Set(context.TODO(), newSessionID.String(), sessionName, expireDuration).Err()
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

	return users, nil
}

func (s *SessionManager) ModifyUser(sessionID uuid.UUID, userModifyID uint64, user *models.UserReq) error {
	marhsalledData, err := json.Marshal(*user)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to marshall user"), err)
	}
	list, err := s.Client.LRange(context.TODO(), sessionID.String(), 0, -1).Result()
	if err != nil {
		return err
	}
	for i, item := range list {
		var user models.UserReq
		err := json.Unmarshal([]byte(item), &user) // Unmarshal Redis list item into struct
		if err != nil {
			return err
		}

		if user.ID == userModifyID {
			var actionStr string
			actionStr, err = s.Client.LSet(context.TODO(), sessionID.String(), int64(i), marhsalledData).Result()
			fmt.Print(actionStr)
			if err != nil {
				return err
			}
			return nil

		}
	}
	return fmt.Errorf("haven't found  users with this ID")
}

func (s *SessionManager) GetUserSessions(userID uint64) ([]Session, error) {

	var err error
	var validSessions []uuid.UUID
	var users []models.UserReq
	var sessions []Session
	var sessionName string
	for _, sessionID := range s.SessionIDs {
		//fmt.Println(sessionID)
		users, err = s.GetUsers(sessionID)
		//fmt.Println(users)
		for _, candidate := range users {
			if candidate.ID == userID {

				if err != nil && err != redis.Nil {
					return nil, err
				}
				if err == nil {
					validSessions = append(validSessions, sessionID)
					/*sessionName, err = s.Client.Get(context.TODO(), sessionID.String()).Result()
					if err != nil && err != redis.Nil {
						return nil, err
					}*/
					sessionName = "someName"

					session := Session{
						SessionID:   sessionID,
						SessionName: sessionName,
						Users:       users,
					}
					sessions = append(sessions, session)
				}
			}
		}
	}
	s.SessionIDs = validSessions

	return sessions, nil
}

//Личный кабинет пользователя
// Хранить в сессии число людей для начала
// Хранить флаг началась сессия/ не началась сессия

package user_repo

import "test_backend_frontend/internal/models"

type IUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	UpdateUserByLogin(login string, user *models.User) error
	DeleteUserByLogin(login string) error
	CreateUser(user *models.User) error
}

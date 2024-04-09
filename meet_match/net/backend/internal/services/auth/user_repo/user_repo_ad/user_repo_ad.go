package repo_adapter

import (
	"test_backend_frontend/internal/models"
	"test_backend_frontend/internal/models/models_da"
	"test_backend_frontend/internal/services/auth/user_repo"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepositoryAdapter struct {
	db *gorm.DB
}

func NewUserRepositoryAdapter(srcDB *gorm.DB) user_repo.IUserRepository {
	return &UserRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *UserRepositoryAdapter) GetUserByID(id uint64) (*models.User, error) {
	var user_da models_da.User
	user_da.ID = id
	tx := repo.db.First(&user_da)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := models_da.FromDaUser(&user_da)
	return &user, nil
}

func (repo *UserRepositoryAdapter) GetUsersByIDs(ids []uint64) ([]models.User, error) {
	var user_das []models_da.User
	tx := repo.db.Where("id IN ?", ids).Find(&user_das)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	var users []models.User

	for _, user_da := range user_das {
		user := models_da.FromDaUser(&user_da)
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepositoryAdapter) GetUserByLogin(login string) (*models.User, error) {
	var user_da models_da.User
	tx := repo.db.Where("login = ?", login).First(&user_da)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := models_da.FromDaUser(&user_da)
	return &user, nil
}

func (repo *UserRepositoryAdapter) UpdateUserByLogin(login string, user *models.User) error {
	user_da := models_da.ToDaUser(*user)
	tx := repo.db.Where("login = ?", login).Updates(user_da)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) DeleteUserByLogin(login string) error {
	tx := repo.db.Where("login = ?", login).Delete(models_da.User{}) // specifically for gorm
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) CreateUser(user *models.User) error {

	tx := repo.db.Create(models_da.ToDaUser(*user))
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating user")
	}
	return nil
}

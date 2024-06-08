package models_da //stands for data_acess

import "test_backend_frontend/internal/models"

type User struct {
	ID       uint64 `gorm:"primaryKey;column:id"`
	Login    string `gorm:"unique;column:login"`
	Password string `gorm:"column:password"`
	Name     string `gorm:"column:name"`
	Surname  string `gorm:"column:surname"`
	Age      int    `gorm:"embedded;column:role"`
	Gender   bool   `gorm:"embedded;column:group"`
}

func FromDaUser(userDa *User) models.User {
	return models.User{
		ID:       userDa.ID,
		Name:     userDa.Name,
		Login:    userDa.Login,
		Password: userDa.Password,
		Surname:  userDa.Surname,
		Age:      userDa.Age,
		Gender:   userDa.Gender,
	}

}

func ToDaUser(user models.User) *User {
	return &User{
		ID:       user.ID,
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Surname:  user.Surname,
		Age:      user.Age,
		Gender:   user.Gender,
	}
}

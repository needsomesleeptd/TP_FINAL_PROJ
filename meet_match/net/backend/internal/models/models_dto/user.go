package models_dto // stands for data_transfer_objects

import (
	"encoding/json"
	"test_backend_frontend/internal/models"
)

type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Age      int    `json:"age"`
	Gender   bool   `json:"gender"`
}

func (u *User) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func FromDtoUser(userDa *User) models.User {
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

func ToDtoUser(user models.User) *User {
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

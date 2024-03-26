package models

type User struct {
	ID       uint64
	Login    string
	Password string
	Name     string
	Surname  string
}

type UserReq struct {
	ID      uint64 `json:"ID"`
	Name    string `json:"Name"`
	Request string `json:"Request"`
	//TODO:: add embeddings
}

func NewUserReq(id uint64, name string, text string) *UserReq {
	req := UserReq{
		ID:      id,
		Name:    name,
		Request: text,
	}
	return &req
}

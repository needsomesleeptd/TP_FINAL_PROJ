package models

type User struct {
	ID       uint64
	Login    string
	Password string
	Name     string
	Surname  string
	Age      int
	Gender   bool // true is a man 0_0, false is women
}

type UserReq struct {
	ID         uint64   `json:"ID" redis:"ID"`
	Name       string   `json:"Name" redis:"Name"`
	Request    string   `json:"Request" redis:"Request"`
	Categories []string `json:"Categories" redis:"Categories"`
}

type SessionStatus int

const (
	Waiting SessionStatus = iota
	Scrolling
	Ended
	ScrollingWasMatched
)

func NewUserReq(id uint64, name string, text string) *UserReq {
	req := UserReq{
		ID:      id,
		Name:    name,
		Request: text,
	}
	return &req
}

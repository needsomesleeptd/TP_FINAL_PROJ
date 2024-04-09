package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Remove card from here to dto + add Id to it

type Card struct {
	ImgUrl   string `json:"image"`
	CardName string `json:"title,card_name"`
	Rating   int    `json:"rating,omitempty"`
}

type ModelRequest struct {
	Query    string `json:"query"`
	Label    string `json:"label"`
	FromLine int    `json:"from_line"`
	ToLine   int    `json:"to_line"`
}

type RecSys struct {
	Url string
}

func New(urlRecSys string) (*RecSys, error) {
	if urlRecSys == "" {
		return &RecSys{}, fmt.Errorf("empty url")
	}

	return &RecSys{Url: urlRecSys}, nil
}

func (r *RecSys) CardsSearch(prompt string, fromLine int, toLine int) ([]Card, error) {
	req := ModelRequest{Query: prompt, Label: "smth", FromLine: fromLine, ToLine: toLine}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		panic(err.Error())
	}

	buffer := bytes.NewBuffer(jsonReq)
	json_resp, err := http.Post(r.Url, "application/json", buffer)
	if err != nil {
		return []Card{}, fmt.Errorf("%s", "Post to model failure")
	}
	fmt.Println(json_resp)
	defer json_resp.Body.Close()

	var arr []Card
	json.NewDecoder(json_resp.Body).Decode(&arr)

	return arr, nil
}

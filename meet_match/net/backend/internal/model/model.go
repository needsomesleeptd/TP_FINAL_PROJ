package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"test_backend_frontend/internal/services/cards/repository"
)

// TODO: Remove card from here to dto + add Id to it

type Card struct {
	Idx      uint64 `json:"idx,omitempty"`
	ImgUrl   string `json:"image"`
	CardName string `json:"title,card_name"`
	Rating   int    `json:"rating,omitempty"`
}

type ModelResponse struct {
	Recs []uint64 `json:"recs"`
}

type ModelRequest struct {
	Query     string `json:"query"`
	SessionID string `json:"session_id"`
	UserID    uint64 `json:"user_id"`
}

type RecSys struct {
	Url     string
	cardRep repository.CardRepository
}

func New(urlRecSys string, cardRepository repository.CardRepository) (*RecSys, error) {
	if urlRecSys == "" {
		return &RecSys{}, fmt.Errorf("empty url")
	}

	return &RecSys{Url: urlRecSys, cardRep: cardRepository}, nil
}

// TODO: refactor
func (r *RecSys) CardsSearch(prompt string, sessionId string, userId uint64) ([]Card, error) {
	req := ModelRequest{
		Query:     prompt,
		SessionID: sessionId,
		UserID:    userId,
	}
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

	var arr ModelResponse
	json.NewDecoder(json_resp.Body).Decode(&arr)

	var cards []Card
	for _, v := range arr.Recs {
		card, err := r.cardRep.GetCard(v)
		if err != nil {
			return []Card{}, fmt.Errorf("%s", "Post to model failure")
		}

		cards = append(cards, Card{
			Idx:      card.Id,
			ImgUrl:   card.ImgUrl,
			CardName: card.CardName,
			Rating:   card.Rating,
		})
	}

	return cards, nil
}

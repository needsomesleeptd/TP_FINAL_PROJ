package rec_model_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"test_backend_frontend/internal/models/models_dto"
	"test_backend_frontend/internal/services/cards/repository"
)

type ModelResponse struct {
	Recs []uint64 `json:"recs"`
}

type ModelRequest struct {
	Query      string   `json:"query"`
	SessionID  string   `json:"session_id"`
	UserID     uint64   `json:"user_id"`
	Categories []string `json:"categories"`
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
func (r *RecSys) CardsSearch(prompt string, sessionId string, userId uint64, categories []string) ([]*models_dto.Card, error) {
	req := ModelRequest{
		Query:      prompt,
		SessionID:  sessionId,
		UserID:     userId,
		Categories: categories,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		panic(err.Error())
	}

	buffer := bytes.NewBuffer(jsonReq)
	json_resp, err := http.Post(r.Url, "application/json", buffer)
	if err != nil {
		return nil, fmt.Errorf("%s-%w", "Post to rec-model-client failure", err)
	}
	defer json_resp.Body.Close()

	var arr ModelResponse
	json.NewDecoder(json_resp.Body).Decode(&arr)

	var cards []*models_dto.Card
	for _, v := range arr.Recs {
		card, err := r.cardRep.GetCard(v)
		if err != nil {
			return nil, fmt.Errorf("%s-%w", "Post to rec-model-client failure", err)
		}

		cards = append(cards, models_dto.ToDTOCard(card))
	}

	return cards, nil
}

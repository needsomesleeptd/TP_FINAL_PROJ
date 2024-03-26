package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func find_csv_row(path string, id int) []string {
	file, err := os.Open("../database/dist.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	targetValue := id

	for _, record := range records {
		val, _ := strconv.Atoi(record[0])
		if val == targetValue {

			return record
		}
	}
	return nil
}

type Card struct {
	Img_url   string `json:"image"`
	Card_name string `json:"title,card_name"`
	Rating    int    `json:"rating,omitempty"`
}

type nn_reqest struct {
	Query string `json:"query"`
	Label string `json:"label"`
}

var tpl = template.Must(template.ParseFiles("../frontend/card.html"))
var card = 0

var array_of_cards []Card


func cards_page(w http.ResponseWriter, r *http.Request) {
	card = 0

	req := r.FormValue("request")
	fmt.Println(req) // В req запрос

	nn_req_struct := nn_reqest{Query: req, Label: "smth"}
	json_req, err := json.Marshal(nn_req_struct)
	if err != nil {
		panic(err.Error())
	}

	buffer := bytes.NewBuffer(json_req)
    url := "http://127.0.0.1:5000/rec"
	json_resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		fmt.Errorf("%s", "Анлак, не получили ответ")
	}
	fmt.Println(json_resp)
	
	defer json_resp.Body.Close()

	json.NewDecoder(json_resp.Body).Decode(&array_of_cards)
	
	

	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards) - 1{
		card += 1
	}
}

func card_page_like(w http.ResponseWriter, r *http.Request) {
	
	fmt.Printf("Карта %d понравилась\n", card)
	
	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards) - 1{
		card += 1
	}
}

func card_page_dislike(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards) - 1{
		card += 1
	}
}

func index_page(w http.ResponseWriter, r *http.Request) {
	var tpl_index = template.Must(template.ParseFiles("../frontend/index.html"))

	tpl_index.Execute(w, nil)

}

func nn_page(w http.ResponseWriter, r *http.Request) {
	//var tpl_index = template.Must(template.ParseFiles("../frontend/index.html"))
	nn_req_struct := nn_reqest{Query: "хочу пива", Label: "smth"}
	json_req, err := json.Marshal(nn_req_struct)
	if err != nil {
		panic(err.Error())
	}

	buffer := bytes.NewBuffer(json_req)
	url := "http://127.0.0.1:5000/rec"
	json_resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		fmt.Errorf("%s", "Анлак, не получили ответ")
	}
	fmt.Println(json_resp)
	var places []map[string]interface{}

	defer json_resp.Body.Close()

	json.NewDecoder(json_resp.Body).Decode(&places)

}

func main() {
	http.HandleFunc("/", index_page)
	http.HandleFunc("/cards/", cards_page)
	http.HandleFunc("/card_like/", card_page_like)
	http.HandleFunc("/card_dislike/", card_page_dislike)
	http.ListenAndServe(":8080", nil)
}

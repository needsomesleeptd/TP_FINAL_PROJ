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

func card_page(w http.ResponseWriter, r *http.Request) {

	//test_card := Card{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg",
	//Card_name: "умный мужик", rating: 5}

	req := r.FormValue("request")

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
	var cards []Card

	defer json_resp.Body.Close()

	json.NewDecoder(json_resp.Body).Decode(&cards)

	tpl.Execute(w, cards[0])

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
	//temp = template.Must(template.ParseGlob("../frontend/*.html"))
	http.HandleFunc("/", index_page)
	http.HandleFunc("/card/", card_page)
	http.HandleFunc("/nn/", nn_page)
	http.ListenAndServe(":8080", nil)
}

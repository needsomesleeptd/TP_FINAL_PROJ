package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Card struct {
	Img_url   string
	Card_name string
	rating    int
}

type query struct {
	qUERY string `json:"Name"`
	Type  string `json:"Type"`
}

var tpl = template.Must(template.ParseFiles("../frontend/card.html"))
var card = 0

var array_of_cards []Card


func cards_page(w http.ResponseWriter, r *http.Request) {
	card = 0
	temp := []Card{{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg",
	Card_name: "умный мужик", rating: 5}, 
	{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg", Card_name: "умный мужик2", rating: 5},
	{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg", Card_name: "умный мужик3", rating: 5},
	{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg",
	Card_name: "умный мужик4", rating: 5}}
	array_of_cards = temp
	
	req := r.FormValue("request")
	fmt.Println(req) // В req запрос

	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards){
		card += 1
	}
}

func card_page_like(w http.ResponseWriter, r *http.Request) {
	
	fmt.Printf("Карта %d понравилась", card)
	
	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards){
		card += 1
	}
}

func card_page_dislike(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards){
		card += 1
	}
}

func index_page(w http.ResponseWriter, r *http.Request) {
	var tpl_index = template.Must(template.ParseFiles("../frontend/index.html"))

	tpl_index.Execute(w, nil)

}

func neural_network(w http.ResponseWriter, r *http.Request) {
	var tpl_index = template.Must(template.ParseFiles("../frontend/index.html"))

	tpl_index.Execute(w, nil)

}

func main() {
	http.HandleFunc("/", index_page)
	http.HandleFunc("/cards/", cards_page)
	http.HandleFunc("/card_like/", card_page_like)
	http.HandleFunc("/card_dislike/", card_page_dislike)
	http.ListenAndServe(":8080", nil)
}

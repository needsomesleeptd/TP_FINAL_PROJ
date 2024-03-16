package main

import (
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

func card_page(w http.ResponseWriter, r *http.Request) {

	test_card := Card{Img_url: "https://media.kudago.com/thumbs/xl/images/list/42/0c/420c5ac9b0836258f52c0b4ee131e6e1.jpg",
		Card_name: "умный мужик", rating: 5}
	tpl.Execute(w, test_card)
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
	http.HandleFunc("/card/", card_page)
	http.ListenAndServe(":8080", nil)
}

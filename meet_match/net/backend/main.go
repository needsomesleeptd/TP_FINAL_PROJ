package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"meetMatch/models"
	sessions "meetMatch/sessions"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TempSessionData struct {
	Session_id   uuid.UUID
	People_count int
}

var sessionManager *sessions.SessionManager

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
	if card < len(array_of_cards)-1 {
		card += 1
	}
}

func card_page_like(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Карта %d понравилась\n", card)

	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards)-1 {
		card += 1
	}
}

func card_page_dislike(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, array_of_cards[card])
	if card < len(array_of_cards)-1 {
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

func session_create_page(w http.ResponseWriter, r *http.Request) {
	req := r.FormValue("request")
	userReq := models.NewUserReq(2, "anyname", req)
	sessionID, err := sessionManager.CreateSession(userReq)
	if err != nil {
		fmt.Errorf(err.Err().Error())
	}

	var tpl_index = template.Must(template.ParseFiles("../frontend/create_session.html"))

	tpl_index.Execute(w, nil)
	s_path := fmt.Sprintf("/session/%u", sessionID)
	http.Redirect(w, r, s_path, http.StatusFound)
}

func session_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sessionID uuid.UUID
	var err error
	sessionID, err = uuid.Parse(vars["id"])
	if err != nil {
		fmt.Print(err)
	}
	users, err := sessionManager.GetUsers(sessionID)
	if err != nil {
		fmt.Print(err)
	}
	var tpl_index = template.Must(template.ParseFiles("../frontend/sessions.html"))
	data := TempSessionData{Session_id: sessionID, People_count: len(users)}
	err = tpl_index.Execute(w, data)
	if err != nil {
		fmt.Print(err)
	}
}

func main() {
	var err error
	sessionManager, err = sessions.NewSessionManager("localhost:6379", "", 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	http.HandleFunc("/", index_page)
	http.HandleFunc("/cards/", cards_page)
	http.HandleFunc("/card_like/", card_page_like)
	http.HandleFunc("/card_dislike/", card_page_dislike)
	http.HandleFunc("/session_create/", session_create_page)
	http.HandleFunc("/session/{id}", session_page)

	http.ListenAndServe(":8080", nil)
}

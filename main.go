package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s !!!", task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var req requestBody                          //переменная для хранекия данных из тела запроса
	decoder := json.NewDecoder(r.Body)           //создаём декодер JSON, привязанный к телу запроса r.Body
	if err := decoder.Decode(&req); err != nil { //декодируем JSON из тела запроса в переменную req, при возникновении ошибки возвращаем код 400
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	task = req.Message
	w.WriteHeader(http.StatusOK) //устанавливаем статус ответа 200
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/task", GetHandler).Methods("GET")
	router.HandleFunc("/api/task", PostHandler).Methods("POST")

	http.ListenAndServe(":8080", router)

}

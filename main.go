package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	if arr := json.NewDecoder(r.Body).Decode(&t); arr != nil {
		http.Error(w, "Invalid recording format", http.StatusBadRequest)
		return
	}
	if err := DB.Create(&t).Error; err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
	}
	//var req requestBody                          //переменная для хранекия данных из тела запроса
	//decoder := json.NewDecoder(r.Body)           //создаём декодер JSON, привязанный к телу запроса r.Body
	//if err := decoder.Decode(&req); err != nil { //декодируем JSON из тела запроса в переменную req, при возникновении ошибки возвращаем код 400
	//	http.Error(w, "Invalid request", http.StatusBadRequest)
	//	return
	//}
	//task = req.Message
	//w.WriteHeader(http.StatusOK) //устанавливаем статус ответа 200
	//fmt.Fprintf(w, "Request: %s sent succesfully", task)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	if err := DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to get tasks list", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Failed to get task", http.StatusNotFound)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Failed to decode update", http.StatusBadRequest)
		return
	}

	allowedFields := map[string]bool{
		"task":    true,
		"is_done": true,
	}
	for key := range allowedFields {
		if !allowedFields[key] {
			delete(updateData, key)
		}
	}
	if err := DB.Model(&task).Updates(updateData).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var task Task

	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Failed to get task", http.StatusNotFound)
		return
	}

	if err := DB.Delete(&task).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func main() {

	InitDB()

	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/message", GetHandler).Methods("GET")            //C
	router.HandleFunc("/api/message", PostHandler).Methods("POST")          //R
	router.HandleFunc("/api/message/{id}", PatchHandler).Methods("PATCH")   //U
	router.HandleFunc("/api/message/{id}", DeleteHandler).Methods("DELETE") //D

	http.ListenAndServe(":8080", router)
}

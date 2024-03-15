package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CraeteUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var userSimpleStorage []User

func main() {
	userSimpleStorage = []User{}

	http.HandleFunc("/items", items)
	http.ListenAndServe(":8090", nil)
}

func handleCreate(createData CraeteUser) User {
	userData := User{
		ID:    uint(time.Now().Unix()),
		Name:  createData.Name,
		Email: createData.Email,
		Age:   createData.Age,
	}

	userSimpleStorage = append(userSimpleStorage, userData)

	return userData
}

func items(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userSimpleStorage)
		return

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var requestData CraeteUser
		err := decoder.Decode(&requestData)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		data := handleCreate(requestData)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return

	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

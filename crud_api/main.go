package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var db *sql.DB

func connectDatabase() {
	connString := "postgresql://tms:codev123@localhost:5432/tms?sslmode=disable"

	var err error
	db, err = sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("error to connect with message %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error when ping with message %s", err.Error())
	}

	log.Println("databae successfully connected")
}

func main() {
	connectDatabase()

	http.HandleFunc("/items", items)

	http.ListenAndServe(":8090", nil)
}

func handleCreate(userRequest UserRequest) (User, error) {
	userData := User{
		ID:    uint(time.Now().Unix()),
		Name:  userRequest.Name,
		Email: userRequest.Email,
		Age:   userRequest.Age,
	}

	_, err := db.ExecContext(context.TODO(), `INSERT INTO "user" ("id", "name", "email", "age") VALUES ($1, $2, $3, $4);`, userData.ID, userData.Name, userData.Email, userData.Age)

	return userData, err
}

func handleUpdate(userRequest UserRequest, id uint) (User, error) {
	userData := User{
		ID:    id,
		Name:  userRequest.Name,
		Email: userRequest.Email,
		Age:   userRequest.Age,
	}

	_, err := db.ExecContext(context.TODO(), `UPDATE "user" SET "name" = $1, "email" = $2, "age" = $3 WHERE "id" = $4;`, userData.Name, userData.Email, userData.Age, userData.ID)

	return userData, err
}

func handleDelete(id uint) error {

	_, err := db.ExecContext(context.TODO(), `DELETE FROM "user" WHERE "id" = $1;`, id)

	return err
}

func handleGetList() ([]User, error) {
	users := []User{}
	rows, err := db.QueryContext(context.TODO(), `SELECT "id", "name", "email", "age" FROM "user"`)
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}

	return users, err
}

func items(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		data, err := handleGetList()
		if err != nil {
			http.Error(w, fmt.Sprintf("error executing query with msg: %s", err.Error()), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(data)
		return

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var requestData UserRequest
		err := decoder.Decode(&requestData)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
		}

		data, err := handleCreate(requestData)
		if err != nil {
			http.Error(w, fmt.Sprintf("error executing query with msg: %s", err.Error()), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return

	case http.MethodPut:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "params id can not be empty", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var requestData UserRequest
		err := decoder.Decode(&requestData)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		idInt, _ := strconv.Atoi(id)
		data, err := handleUpdate(requestData, uint(idInt))
		if err != nil {
			http.Error(w, fmt.Sprintf("error executing query with msg: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "params id can not be empty", http.StatusBadRequest)
			return
		}

		idInt, _ := strconv.Atoi(id)
		err := handleDelete(uint(idInt))
		if err != nil {
			http.Error(w, fmt.Sprintf("error executing query with msg: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

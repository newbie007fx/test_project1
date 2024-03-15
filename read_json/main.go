package main

import (
	"encoding/json"
	"log"
	"os"
)

type DataUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	log.Println("updating file json")
	dataBytes, err := os.ReadFile("./file.json")
	if err != nil {
		log.Println(err.Error())
		return
	}

	user := new(DataUser)
	err = json.Unmarshal(dataBytes, user)
	if err != nil {
		log.Println(err.Error())
		return
	}

	user.Age += 1
	user.Email = "johndoe@example.com"

	dataBytes, err = json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = os.WriteFile("./file.json", dataBytes, 0644)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("success update file json")
}

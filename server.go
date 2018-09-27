package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type User struct{ Name string }

var users = make(map[int]User)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[7:]
	index, err := strconv.Atoi(key)
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("Invalid user id: %s", key), http.StatusNotFound)
		return
	}
	user, ok := users[index]
	if ok {
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, fmt.Sprintf("Invalid user id: %s", key), http.StatusNotFound)
	}
}

func main() {
	users[0] = User{"admin"}
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	log.Print("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

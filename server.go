package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var users = make(map[int]User)

var lock sync.Mutex

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		lock.Lock()
		json.NewEncoder(w).Encode(users)
		lock.Unlock()
	} else if r.Method == http.MethodPost {
		var user User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Error during read body", http.StatusBadRequest)
			return
		}
		if err := r.Body.Close(); err != nil {
			log.Printf("Error during closing request: %v", err)
			http.Error(w, "Error during read body", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &user); err != nil {
			log.Printf("Error during parsing request: %v", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			http.Error(w, "Error during parsing body", http.StatusBadRequest)
			return
		}
		lock.Lock()
		index := len(users)
		users[index] = user
		lock.Unlock()
	} else {
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[7:]
	index, err := strconv.Atoi(key)
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("Invalid user id: %s", key), http.StatusNotFound)
		return
	}
	lock.Lock()
	user, ok := users[index]
	lock.Unlock()
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

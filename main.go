package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id      int
	Name    string
	Balance float32
}

// type Quest struct {
// 	id   int
// 	name string
// 	cost float32
// }

var db *sql.DB

func main() {
	connStr := "user=username1 password=password1 dbname=db host=db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	// router.HandleFunc("/quests", createQuest).Methods("POST")
	// router.HandleFunc("/complete", completeQuest).Methods("POST")
	// router.HandleFunc("/history/{user_id}", getUserHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Print(user.Name)

	query := "INSERT INTO users(name, balance) VALUES($1, $2) RETURNING id"
	err = db.QueryRow(query, user.Name, user.Balance).Scan(&user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

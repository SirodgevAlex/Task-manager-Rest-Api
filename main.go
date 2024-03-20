package main

import (
	"database/sql"
	"encoding/json"
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
	var err error
	connStr := "user=postgres password=1234 dbname=db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
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

	query := "INSERT INTO users(name, balance) VALUES($1, $2) RETURNING id"
	
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() 
	
	err = tx.QueryRow(query, user.Name, user.Balance).Scan(&user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
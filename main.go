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

type Quest struct {
	Id   int
	Name string
	Cost float32
}

type Task struct {
	Id      int
	UserId  int
	QuestId int
}

type UserHistory struct {
	completedQuests []Quest
	Balance         float32
}

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
	router.HandleFunc("/quests", createQuest).Methods("POST")
	router.HandleFunc("/complete", completeQuest).Methods("POST")
	router.HandleFunc("/history/{userId}", getUserHistory).Methods("GET")

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

func createQuest(w http.ResponseWriter, r *http.Request) {
	var quest Quest
	err := json.NewDecoder(r.Body).Decode(&quest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO quests(name, cost) VALUES($1, $2) RETURNING id"

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	err = tx.QueryRow(query, quest.Name, quest.Cost).Scan(&quest.Id)
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
	json.NewEncoder(w).Encode(quest)
}

func completeQuest(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM completedTasks WHERE userId = $1 AND questId = $2;", task.UserId, task.QuestId).Scan(&count)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		tx.Rollback()
		http.Error(w, "Task already completed by user", http.StatusBadRequest)
		return
	}

	var cost float32
	err = tx.QueryRow("SELECT cost FROM quests WHERE id = $1;", task.QuestId).Scan(&cost)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2;", cost, task.UserId)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO completedTasks (userId, questId) VALUES ($1, $2);", task.UserId, task.QuestId)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	var userHistory UserHistory
	userId := mux.Vars(r)["userId"]

	rows, err := db.Query("SELECT * FROM quests q JOIN completedTasks ct ON q.Id = ct.questId WHERE ct.userId = $1", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var quest Quest
		var task  Task
		if err := rows.Scan(&quest.Id, &quest.Name, &quest.Cost, &task.Id, &task.UserId, &task.QuestId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userHistory.completedQuests = append(userHistory.completedQuests, quest)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var balance float32
	err = db.QueryRow("SELECT balance from users where Id = $1", userId).Scan(&balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userHistory.Balance = balance

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userHistory)
}

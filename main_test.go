package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	reqBody := []byte(`{"Name": "Test User", "Balance": 100.5}`)
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()

	createUser(responseWriter, req)

	if status := responseWriter.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"Id":1,"Name":"Test User","Balance":100.5}`

	if responseWriter.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseWriter.Body.String(), expected)
	}
}

func TestCreateQuest(t *testing.T) {
	reqBody := []byte(`{"Name": "Test User", "Cost": 10.5}`)
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/quests", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	responseWriter := httptest.NewRecorder()

	createUser(responseWriter, req)

	if status := responseWriter.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"Id":1,"Name":"Test User","Cost":10.5}`

	if responseWriter.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseWriter.Body.String(), expected)
	}
}

func TestMain(m *testing.M) {
	TestCreateUser(&testing.T{})
	TestCreateQuest(&testing.T{})
}

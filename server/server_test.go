package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/VladislavLisovenko/task_management/client/entities"
	"github.com/VladislavLisovenko/task_management/server/handlers"
)

const (
	userName = "Bob"
)

func getUserByName(name string) entities.User {
	user := entities.User{Name: name}

	encodedMessage, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", bytes.NewReader(encodedMessage))
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UserByName)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		log.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("id") == "" {
		log.Fatal("response must contain id")
	}

	id, err := strconv.Atoi(rr.Header().Get("id"))
	if err != nil {
		log.Fatal("id is not a number")
	}

	user.SetID(id)

	return user
}

func tackList(tlf entities.TaskListFilter) []entities.Task {
	encodedMessage, err := json.Marshal(tlf)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(encodedMessage))
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.TaskList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		log.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var tasks []entities.Task
	err = json.NewDecoder(rr.Body).Decode(&tasks)
	if err != nil {
		log.Fatal(err.Error())
	}

	return tasks
}

func TestMainUsersGet(t *testing.T) {
	user := getUserByName(userName)
	if user.GetID() == 0 {
		t.Fatal("error user creating")
	}
}

func TestMainTasksPost(t *testing.T) {
	user := getUserByName(userName)
	task := entities.Task{
		Description:    "Some task",
		ExpirationDate: time.Date(2024, 2, 25, 20, 0, 0, 0, time.Local),
		Done:           false,
		User:           user,
	}

	encodedMessage, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(encodedMessage))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("id") == "" {
		t.Error("response must contain id")
	}
}

func TestTasksGet(t *testing.T) {
	user := getUserByName(userName)
	tlf := entities.TaskListFilter{User: user}
	taskList := tackList(tlf)
	if len(taskList) == 0 {
		t.Fatal("task list should not be empty")
	}
}

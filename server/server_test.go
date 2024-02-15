package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/VladislavLisovenko/task_management/client/entities"
	"github.com/VladislavLisovenko/task_management/server/handlers"

	"github.com/stretchr/testify/require"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	userName = "Andrey"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func getUserByName(name string) (entities.User, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/users?name="+userName, nil)
	if err != nil {
		return entities.User{}, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UserByName)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		return entities.User{}, errors.New("handler returned wrong status code")
	}

	if rr.Header().Get("id") == "" {
		return entities.User{}, errors.New("response must contain id")
	}

	id, err := strconv.Atoi(rr.Header().Get("id"))
	if err != nil {
		return entities.User{}, errors.New("id is not a number")
	}

	user := entities.User{Name: userName}
	user.SetID(id)

	return user, nil
}

func createUser(name string) (entities.User, error) {
	user := entities.User{Name: name}

	encodedMessage, err := json.Marshal(user)
	if err != nil {
		return entities.User{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, "/", bytes.NewReader(encodedMessage))
	if err != nil {
		return entities.User{}, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		return entities.User{}, errors.New(rr.Header().Get("error"))
	}

	id, err := strconv.Atoi(rr.Header().Get("id"))
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{Name: name, ID: id}, nil
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
	user, err := getUserByName(userName)
	require.NoError(t, err)
	require.NotEqual(t, 0, user.GetID())
}

func TestMainUsersPut(t *testing.T) {
	_, err := createUser(userName)
	require.EqualError(t, err, "такой пользователь уже существует")

	user, err := createUser(randSeq(10))
	require.Equal(t, nil, err)
	require.NotEqual(t, 0, user.GetID())
}

func TestMainTasksPost(t *testing.T) {
	user, err := getUserByName(userName)
	require.NoError(t, err)
	task := entities.Task{
		Description:    "Some task",
		ExpirationDate: time.Date(2024, 2, 25, 20, 0, 0, 0, time.Local),
		Done:           false,
		User:           user,
	}

	encodedMessage, err := json.Marshal(task)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(encodedMessage))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		require.Equal(t, http.StatusOK, status)
	}

	require.NotEmpty(t, rr.Header().Get("id"))
}

func TestTasksGet(t *testing.T) {
	user, err := getUserByName(userName)
	require.NoError(t, err)
	tlf := entities.TaskListFilter{User: user}
	taskList := tackList(tlf)
	require.NotEmpty(t, taskList)
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

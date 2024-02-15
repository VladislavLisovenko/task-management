package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/server/db"
	"github.com/VladislavLisovenko/task_management/server/entities"
)

func taskListfilterIsValid(tlf entities.TaskListFilter) bool {
	return tlf.ExpirationDateTo.After(tlf.ExpirationDateFrom) ||

		tlf.ExpirationDateTo.Equal(tlf.ExpirationDateFrom)
}

func TaskList(w http.ResponseWriter, r *http.Request) {
	tlf := decodeEntity[entities.TaskListFilter](w, r)

	if !taskListfilterIsValid(tlf) {
		w.Header().Add("error", "Неверно задан период отбора")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	taskList, err := db.TaskList(tlf)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	taskListDecoded, err := json.Marshal(taskList)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = w.Write(taskListDecoded)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := decodeEntity[entities.Task](w, r)

	err := db.UpdateTask(task)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	task := decodeEntity[entities.Task](w, r)

	if task.CreationDate.After(task.ExpirationDate) {
		w.Header().Add("error", "Неверно указан срок выполнения, он должен быть позже даты создания")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	taskID, err := db.AddTask(task)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.Header().Add("id", strconv.Itoa(taskID))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := decodeEntityID(w, r)

	err := db.DeleteTask(taskID)
	if err != nil {
		w.Header().Add("error", err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}
}

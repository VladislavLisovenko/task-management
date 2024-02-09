package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/server/db"
	"github.com/VladislavLisovenko/task_management/server/entities"
)

func TaskList(w http.ResponseWriter, r *http.Request) {
	tlf := decodeEntity[entities.TaskListFilter](w, r)
	taskList, err := db.TaskList(tlf)
	if err != nil {
		fmt.Println(err.Error())
	}
	taskListDecoded, err := json.Marshal(taskList)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = w.Write(taskListDecoded)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := decodeEntity[entities.Task](w, r)

	err := db.UpdateTask(task)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	task := decodeEntity[entities.Task](w, r)

	taskID, err := db.AddTask(task)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Add("id", strconv.Itoa(taskID))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := decodeEntityID(w, r)

	err := db.DeleteTask(taskID)
	if err != nil {
		fmt.Println(err.Error())
	}
}

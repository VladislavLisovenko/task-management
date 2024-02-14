package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/server/db"
	"github.com/VladislavLisovenko/task_management/server/entities"
	"github.com/go-chi/chi"
)

func UserByName(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "name")

	if userName == "" {
		w.Header().Add("error", "Имя пользователя не указано")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := db.UserByName(userName)
	if err != nil {
		w.Header().Add("error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("id", strconv.Itoa(user.GetID()))
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userDecoded := decodeEntity[entities.User](w, r)

	user, err := db.AddUser(userDecoded.Name)
	if err != nil {
		w.Header().Add("error", fmt.Sprint(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("id", strconv.Itoa(user.GetID()))
}

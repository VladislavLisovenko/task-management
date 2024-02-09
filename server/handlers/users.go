package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/server/db"
	"github.com/VladislavLisovenko/task_management/server/entities"
)

func UserByName(w http.ResponseWriter, r *http.Request) {
	userDecoded := decodeEntity[entities.User](w, r)

	user, err := db.UserByName(userDecoded.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Add("id", strconv.Itoa(user.GetID()))
}

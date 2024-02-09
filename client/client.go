package main

import (
	"fmt"

	"github.com/VladislavLisovenko/task_management/client/menu"
)

func main() {
	url := "http://localhost:8080"

	user, err := menu.LogIn(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	quit := false
	for !quit {
		menu.ShowMenu()

		var response string
		fmt.Scanln(&response)

		switch response {
		case "q":
			quit = true
		case "1":
			menu.AddTask(url, user)
		case "2":
			menu.RemoveTask(url, user)
		case "3":
			menu.EditTask(url, user)
		case "4":
			menu.ListTask(url, user)
		case "5":
			menu.ListTaskWithFilter(url, user)
		default:
			fmt.Println("Номер команды введён неверно, попробуйте ещё раз.")
		}
	}
}

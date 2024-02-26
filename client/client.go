package main

import (
	"fmt"

	"github.com/VladislavLisovenko/task_management/client/entities"
	"github.com/VladislavLisovenko/task_management/client/menu"
)

func main() {
	url := "http://localhost:8080"
	user := entities.User{}
	quit := false

	for !quit {
		menu.ShowLoginMenu()

		var response string

		fmt.Scanln(&response)

		switch response {
		case "q":
			return
		case "1":
			usr, err := menu.SignIn(url)
			if err != nil {
				fmt.Println(err)
				continue
			}

			user = usr
			quit = true
		case "2":
			usr, err := menu.LogIn(url)
			if err != nil {
				fmt.Println(err)
				continue
			}
			user = usr
			quit = true
		default:

			fmt.Println("Номер команды введён неверно, попробуйте ещё раз.")
		}

		if user.GetID() == 0 {
			fmt.Println("Необходимо авторизоваться")
		}
	}

	for {
		menu.ShowMainMenu()
		var response string
		fmt.Scanln(&response)
		switch response {
		case "q":
			return
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

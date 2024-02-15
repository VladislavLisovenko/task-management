package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladislavLisovenko/task_management/server/handlers"
	"github.com/go-chi/chi"
)

func main() {
	fmt.Println("Server started...")

	url := ""

	port := "8080"

	router := chi.NewRouter()

	router.Get("/users/{name:[a-zA-Zа-яА-Я0-9]+}", handlers.UserByName)

	router.Put("/users", handlers.AddUser)

	router.Put("/tasks", handlers.AddTask)

	router.Delete("/tasks/{id:[0-9]+}", handlers.DeleteTask)

	router.Post("/tasks/{id:[0-9]+}", handlers.UpdateTask)

	router.Post("/tasks", handlers.TaskList)

	srv := &http.Server{
		ReadTimeout: 5 * time.Second,

		WriteTimeout: 10 * time.Second,

		Handler: router,
	}

	srv.Addr = fmt.Sprintf("%s:%s", url, port)

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

package httprequests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/client/entities"
)

func SendRequest(url string, method string, message []byte) (*http.Response, error) {
	getRequest, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewReader(message))
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	response, err := httpClient.Do(getRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ProcessEntity[T entities.HasID](url string, httpMethod string, entity T) {
	encMessage, err := json.Marshal(entity)
	if err != nil {
		log.Fatalln(err.Error())
	}

	response, err := SendRequest(url, httpMethod, encMessage)
	if err != nil {
		log.Fatalf("%T %s error: %s", entity, httpMethod, err.Error())
	}
	defer response.Body.Close()

	headerID := response.Header.Get("id")
	if headerID != "" {
		id, err1 := strconv.Atoi(headerID)
		if err1 != nil {
			fmt.Println(err1.Error())
		}
		entity.SetID(id)
	}
}

func ProcessTaskList(url string, tlf entities.TaskListFilter) []entities.Task {
	msgEncoded, err := json.Marshal(tlf)
	if err != nil {
		log.Fatalf("%s on GET error: %s", url, err.Error())
	}
	response, err := SendRequest(url, http.MethodGet, msgEncoded)
	if err != nil {
		log.Fatalf("%s on GET error: %s", url, err.Error())
	}
	defer response.Body.Close()

	var tasks []entities.Task
	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		fmt.Printf("%s on GET error: %s", url, err.Error())
	}

	return tasks
}

func ProcessEntityRemoving(url string) {
	response, err := SendRequest(url, http.MethodDelete, nil)
	if err != nil {
		log.Fatalf("%s on DELETE error: %s", url, err.Error())
	}
	defer response.Body.Close()
}

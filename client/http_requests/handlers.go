package httprequests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladislavLisovenko/task_management/client/entities"
)

func SendGetRequest(url string) (*http.Response, error) {
	getRequest, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
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

func ProcessEntity[T entities.HasID](url string, httpMethod string, entity T) error {
	encMessage, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	response, err := SendRequest(url, httpMethod, encMessage)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New(response.Header.Get("error"))
	}

	headerID := response.Header.Get("id")
	if headerID != "" {
		id, err1 := strconv.Atoi(headerID)
		if err1 != nil {
			return err1
		}
		entity.SetID(id)
	}

	return nil
}

func ProcessTaskList(url string, tlf entities.TaskListFilter) ([]entities.Task, error) {
	msgEncoded, err := json.Marshal(tlf)
	if err != nil {
		return []entities.Task{}, err
	}
	response, err := SendRequest(url, http.MethodPost, msgEncoded)
	if err != nil {
		return []entities.Task{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []entities.Task{}, errors.New(response.Header.Get("error"))
	}

	var tasks []entities.Task
	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return []entities.Task{}, err
	}

	return tasks, nil
}

func ProcessEntityRemoving(url string) error {
	response, err := SendRequest(url, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println(response.Header.Get("error"))
		return errors.New(response.Header.Get("error"))
	}

	return nil
}

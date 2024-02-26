package db

import (
	"errors"

	"github.com/VladislavLisovenko/task_management/server/entities"
)

func UserByName(userName string) (entities.User, error) {
	queryString := `
	SELECT 
		ID
	FROM 
		PUBLIC.USERS
	WHERE 
		NAME = $1`
	row := database.QueryRow(queryString, userName)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{ID: userID, Name: userName}, nil
}

func AddUser(userName string) (entities.User, error) {
	user, _ := UserByName(userName)
	if user.GetID() != 0 {
		return entities.User{}, errors.New("такой пользователь уже существует")
	}
	lastInsertedID := 0
	queryString := `
	INSERT INTO PUBLIC.USERS 
		(NAME)
	VALUES 
		($1) 
	RETURNING ID`

	row := database.QueryRow(queryString, userName)
	err := row.Scan(&lastInsertedID)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{ID: lastInsertedID, Name: userName}, nil
}

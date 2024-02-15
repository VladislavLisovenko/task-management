package db

import (
	"strconv"
	"strings"
	"time"

	"github.com/VladislavLisovenko/task_management/server/entities"
)

func AddTask(task entities.Task) (int, error) {
	description := task.Description

	expirationDate := task.ExpirationDate

	done := task.Done

	user := task.User

	taskID := 0

	queryString := `

	INSERT INTO PUBLIC.TASKS (DESCRIPTION, EXPIRATION_DATE, USER_ID, DONE)

	VALUES ($1, $2, $3, $4) 

	RETURNING ID`

	row := database.QueryRow(queryString, description, expirationDate, user.GetID(), done)

	err := row.Scan(&taskID)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func UpdateTask(task entities.Task) error {
	taskID := task.GetID()

	description := task.Description

	expirationDate := task.ExpirationDate

	userID := task.User.GetID()

	done := task.Done

	queryString := `

	UPDATE PUBLIC.tasks

	SET 

		DESCRIPTION = $1,

		EXPIRATION_DATE = $2,

		USER_ID = $3,

		DONE = $4

	WHERE 

		ID = $5`

	_, err := database.Exec(queryString, description, expirationDate, userID, done, taskID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(taskID int) error {
	queryString := `DELETE FROM PUBLIC.tasks WHERE ID = $1`

	_, err := database.Exec(queryString, taskID)
	if err != nil {
		return err
	}

	return nil
}

func taskListQuery(tlf entities.TaskListFilter) (string, []any) {
	params := make([]any, 0)

	params = append(params, tlf.User.GetID())

	sb := &strings.Builder{}

	sb.WriteString(`

	SELECT 

		T.ID AS TASK_ID,

		T.DESCRIPTION AS DESCRIPTION,

		T.EXPIRATION_DATE AS EXPIRATION_DATE,

		T.USER_ID AS USER_ID,

		U.NAME AS USER_NAME,

		T.DONE AS DONE

	FROM 

		PUBLIC.TASKS AS T

	LEFT JOIN PUBLIC.USERS AS U 

		ON T.USER_ID = U.ID

	WHERE USER_ID = $1`)

	if !tlf.FilterIsSet {
		return sb.String(), params
	}

	if tlf.Description != "" {
		sb.WriteString(` AND T.DESCRIPTION LIKE '%` + tlf.Description + `%'`)
	}

	if !tlf.ExpirationDateFrom.IsZero() {
		params = append(params, tlf.ExpirationDateFrom)

		str := ` AND T.EXPIRATION_DATE >= $` + strconv.Itoa(len(params))

		sb.WriteString(str)
	}

	if !tlf.ExpirationDateTo.IsZero() {
		params = append(params, tlf.ExpirationDateTo)

		str := ` AND T.EXPIRATION_DATE <= $` + strconv.Itoa(len(params))

		sb.WriteString(str)
	}

	if tlf.ConsiderDone {
		params = append(params, tlf.Done)

		str := ` AND T.DONE = $` + strconv.Itoa(len(params))

		sb.WriteString(str)
	}

	return sb.String(), params
}

func TaskList(tlf entities.TaskListFilter) ([]entities.Task, error) {
	query, params := taskListQuery(tlf)

	taskRows, err := database.Query(query, params...)
	if err != nil {
		return nil, err
	}

	defer taskRows.Close()

	if taskRows.Err() != nil {
		return nil, taskRows.Err()
	}

	taskList := make([]entities.Task, 0)

	for taskRows.Next() {
		var taskID int

		var description string

		var expirationDate time.Time

		var userID int

		var userName string

		var done bool

		if err = taskRows.Scan(

			&taskID,

			&description,

			&expirationDate,

			&userID,

			&userName,

			&done); err != nil {
			return nil, err
		}

		task := entities.Task{
			ID: taskID,

			Description: description,

			ExpirationDate: expirationDate,

			User: entities.User{ID: userID, Name: userName},

			Done: done,
		}

		taskList = append(taskList, task)
	}

	return taskList, nil
}

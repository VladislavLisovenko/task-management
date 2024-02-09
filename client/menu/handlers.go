package menu

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VladislavLisovenko/task_management/client/entities"
	httprequests "github.com/VladislavLisovenko/task_management/client/http_requests"
)

const (
	timeLayout = "02/01/2006"
	dateFormat = "dd/MM/yyyy"
	yes        = "да"
	no         = "нет"
)

func formatBool(b bool) string {
	if b {
		return yes
	}
	return no
}

func StringAsDate(str string) (time.Time, error) {
	date, err := time.Parse(timeLayout, str)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func scanDate(hint string) (time.Time, error) {
	fmt.Println(hint, fmt.Sprintf("(в формате %s):", dateFormat))
	var str string
	fmt.Scanln(&str)
	if str == "" {
		return time.Time{}, nil
	}

	return StringAsDate(str)
}

func ShowMenu() {
	fmt.Println("Введите номер команды 1-5 или q для выхода:")
	fmt.Println("1. Добавить задачу")
	fmt.Println("2. Удалить задачу")
	fmt.Println("3. Изменить задачу")
	fmt.Println("4. Показать список всех задач")
	fmt.Println("5. Показать список задач с отбором")
}

func LogIn(url string) (entities.User, error) {
	fmt.Println("Введите имя пользователя или 'q' для выхода:")
	var userName string
	fmt.Scanln(&userName)
	if strings.ToLower(userName) == "q" {
		return entities.User{}, errors.New("завершение работы пользователем")
	}

	addr := fmt.Sprintf("%s/users", url)
	user := entities.User{Name: userName}
	httprequests.ProcessEntity[*entities.User](addr, http.MethodGet, &user)
	if user.GetID() == 0 {
		return entities.User{}, errors.New("не удалось получить информацию о пользователе, убедитесь что сервер доступен")
	}

	return user, nil
}

func ShowTaskList(tasks []entities.Task) {
	fmt.Println()
	fmt.Println("Список задач пользователя:")
	fmt.Println("ID\tСрок выполнениия\tВыполнена\tОписание задачи")
	for _, t := range tasks {
		fmt.Printf("%d\t%s\t\t%s\t\t%s\n",
			t.GetID(),
			t.ExpirationDate.Local().Format(timeLayout),
			formatBool(t.Done),
			t.Description)
	}
	fmt.Println()
}

func AddTask(url string, user entities.User) {
	fmt.Println("Введите описание задачи:")
	var description string
	fmt.Scanln(&description)

	expirationDate, err := scanDate("Введите крайний срок выполнения")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if expirationDate.IsZero() {
		fmt.Println("Срок выполнения не может быть пустым")
		return
	}

	task := &entities.Task{Description: description, ExpirationDate: expirationDate, User: user}
	addr := fmt.Sprintf("%s/tasks", url)
	httprequests.ProcessEntity[*entities.Task](addr, http.MethodPost, task)
}

func RemoveTask(url string, user entities.User) {
	addr := fmt.Sprintf("%s/tasks", url)
	tasks := httprequests.ProcessTaskList(addr, entities.TaskListFilter{User: user})
	ShowTaskList(tasks)

	fmt.Println("Введите ID задачи:")
	var taskID string
	fmt.Scanln(&taskID)
	taskIDIncorrect := true
	for _, t := range tasks {
		if taskID == strconv.Itoa(t.GetID()) {
			taskIDIncorrect = false
			break
		}
	}
	if taskIDIncorrect {
		fmt.Println("ID задачи не найден.")
		return
	}
	addr = fmt.Sprintf("%s/tasks/%s", url, taskID)
	httprequests.ProcessEntityRemoving(addr)
}

func EditTask(url string, user entities.User) {
	addr := fmt.Sprintf("%s/tasks", url)
	tasks := httprequests.ProcessTaskList(addr, entities.TaskListFilter{User: user})

	fmt.Println("")
	fmt.Println("Введите ID задачи:")
	var taskID string
	fmt.Scanln(&taskID)
	taskIDIncorrect := true
	var task entities.Task
	for _, t := range tasks {
		if taskID == strconv.Itoa(t.GetID()) {
			task = t
			taskIDIncorrect = false
			break
		}
	}
	if taskIDIncorrect {
		fmt.Println("ID задачи не найден.")
		return
	}

	fmt.Println("Введите новые значения только для тех атрибутов, которые необходимо изменить")
	taskChanged := false

	fmt.Printf("Введите описание задачи\nТекущее значение: %s\n", task.Description)
	var description string
	fmt.Scanln(&description)
	if description != "" {
		task.Description = description
		taskChanged = true
	}

	hint := fmt.Sprintf("Введите крайний срок выполнения\nТекущее значение: %s\n", task.ExpirationDate.Format(timeLayout))
	expirationDate, err := scanDate(hint)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !expirationDate.IsZero() {
		task.ExpirationDate = expirationDate
		taskChanged = true
	}

	fmt.Println("Введите признак выполнения (возможные значения - да/нет):")
	fmt.Printf("Текущее значение: %s\n", formatBool(task.Done))
	var doneStr string
	fmt.Scanln(&doneStr)
	if doneStr != "" {
		switch strings.ToLower(doneStr) {
		case yes:
			task.Done = true
			taskChanged = true
		case no:
			task.Done = false
			taskChanged = true
		default:
			fmt.Println("Введено неверное значение, признак выполнения оставлен без изменений")
		}
	}

	if taskChanged {
		addr = fmt.Sprintf("%s/tasks/%d", url, task.GetID())
		httprequests.ProcessEntity[*entities.Task](addr, http.MethodPost, &task)
	}
}

func ListTask(url string, user entities.User) {
	addr := fmt.Sprintf("%s/tasks", url)
	tasks := httprequests.ProcessTaskList(addr, entities.TaskListFilter{User: user})
	ShowTaskList(tasks)
}

func ListTaskWithFilter(url string, user entities.User) {
	fmt.Println("Введите значения полей отбора (оставьте поле пустым, если отбор по нему не нужен) и нажмите Enter:")
	taskListFilter := &entities.TaskListFilter{User: user}

	fmt.Println("Значение отбора по полю 'Описание':")
	var description string
	fmt.Scanln(&description)
	if description != "" {
		taskListFilter.SetDescription(description)
	}

	fmt.Println("Значение отбора по полю 'Дата выполнения' устанавливается в виде двух дат - начала и конца периода.")

	startDate, err := scanDate("Начало периода")
	if err != nil {
		fmt.Println("Не удалось определить дату:", err.Error())
	}
	if !startDate.IsZero() {
		taskListFilter.SetExpirationDateFrom(startDate)
	}

	endDate, err := scanDate("Конец периода")
	if err != nil {
		fmt.Println("Не удалось определить дату:", err.Error())
	}
	if !endDate.IsZero() {
		taskListFilter.SetExpirationDateTo(endDate)
	}

	fmt.Println("Значение отбора по полю 'Выполнена' (возможные значения - да/нет):")
	var doneStr string
	fmt.Scanln(&doneStr)
	if doneStr != "" {
		if strings.ToLower(doneStr) == yes {
			taskListFilter.SetDone(true)
		} else if strings.ToLower(doneStr) == no {
			taskListFilter.SetDone(false)
		}
	}

	addr := fmt.Sprintf("%s/tasks", url)
	tasks := httprequests.ProcessTaskList(addr, *taskListFilter)
	ShowTaskList(tasks)
}

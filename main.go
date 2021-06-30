package main

import (
	"30/storage"
	"fmt"
	"strconv"
)

func main() {

	connstr := "postgres://postgres:ts950sdx@localhost/tasks"

	db, err := storage.New(connstr)

	if err != nil {
		fmt.Println(err)
		return
	}
	// Создание новой записи в БД
	var task = storage.Task{
		ID:         1,
		Opened:     2,
		Closed:     1,
		AuthorID:   1,
		AssignedID: 1,
		Title:      "В круге первом",
		Content:    "QQQQQ",
	}

	new_id, err := db.NewTask(task)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(task, db)
	fmt.Println("Создание записи в БД, id=", new_id)

	// Вывод массива данных из БД
	var tasks []storage.Task
	for i := 1; i < 10; i++ {
		t := storage.Task{Title: strconv.Itoa(i), Content: strconv.Itoa(i)}
		tasks = append(tasks, t)
	}

	fmt.Println("Создание массива данных, tasks=", tasks)

	//Удаление записи в БД

	err = db.Delete(4)
	if err != nil {
		fmt.Println(err)
		return
	}
	upd := storage.Task{
		ID:         5,
		Opened:     2,
		Closed:     1,
		AuthorID:   1,
		AssignedID: 1,
		Title:      "first time",
		Content:    "QQQQQ",
	}
	err = db.Update(upd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Обновление записи БД ", upd)

	//Получение списка задач по метке

	label, err := db.Label("see glass")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Список задач по метке ", label)

	//Получение списка задач по автору

	autor, err := db.Tasks(0, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Список задач по автору ", autor)
}

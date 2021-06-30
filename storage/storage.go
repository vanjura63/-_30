package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content,opened,closed)
		VALUES ($1, $2,$3,$4) RETURNING id;
		`,
		t.Title,
		t.Content,
		t.Opened,
		t.Closed,
	).Scan(&id)
	return id, err
}

// Удаление записи из БД по ID
func (s *Storage) Delete(t int) error {

	_, err := s.db.Exec(context.Background(), `
		DELETE  from tasks where id=$1
		
		`,
		t,
	)
	return err
}

// Изменение записи БД по ID
func (s *Storage) Update(t Task) error {

	_, err := s.db.Exec(context.Background(), `
		UPDATE  tasks SET opened=$1,closed=$2,author_ID=$3,assigned_ID=$4,title=$5,content=$6
		WHERE id=$7
		`,
		t.Opened, t.Closed, t.AuthorID, t.AssignedID, t.Title, t.Content, t.ID)

	return err
}

//Получение списка задач по метке
func (s *Storage) Label(label_name string) ([]Task, error) {

	rows, err := s.db.Query(context.Background(), `
	SELECT tasks.id,Opened,Closed,Author_ID,Assigned_ID,Title,Content FROM tasks JOIN tasks_labels ON tasks_labels.task_id=tasks.id JOIN labels ON labels.id=tasks_labels.label_id
	WHERE labels.name=$1
	`,
		label_name,
	)

	if err != nil {
		return nil, err
	}
	var tasks []Task

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

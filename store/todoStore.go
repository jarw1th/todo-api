package store

import (
	models "ToDoProject/models"
	"ToDoProject/utils"
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type TodoStore struct {
	DB *sql.DB
}

func NewTodoStore(connStr string) (*TodoStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &TodoStore{DB: db}, nil
}

func (s *TodoStore) Create(userId int, title string, description string) (models.Todo, error) {
	var t models.Todo
	err := s.DB.QueryRow(
		"INSERT INTO todos(user_id, title, description, done) VALUES($1, $2, $3, $4) RETURNING id, user_id, title, description, created_at, done",
		userId, title, description, false,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)
	if err == nil {
		s.recordHistory(t.ID, userId, models.Todo{}, t)
	}
	return t, err
}

func (s *TodoStore) List(userId int) ([]models.Todo, error) {
	rows, err := s.DB.Query("SELECT id, user_id, title, description, created_at, done FROM todos WHERE user_id=$1 ORDER BY id", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		rows.Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)
		todos = append(todos, t)
	}
	return todos, nil
}

func (s *TodoStore) FilteredList(userId int, m models.TodoQueries) ([]models.Todo, error) {
	query := "SELECT id, user_id, title, description, created_at, done FROM todos"
	var args []interface{}
	var conditions []string

	if m.Done != nil {
		args = append(args, *m.Done)
		conditions = append(conditions, "done = $"+strconv.Itoa(len(args)))
	}

	if m.Title != nil && *m.Title != "" {
		args = append(args, "%"+*m.Title+"%")
		conditions = append(conditions, "description ILIKE $"+strconv.Itoa(len(args)))
	}

	if m.Description != nil && *m.Description != "" {
		args = append(args, "%"+*m.Description+"%")
		conditions = append(conditions, "title ILIKE $"+strconv.Itoa(len(args)))
	}

	if m.Timestamp != nil && *m.Timestamp != "" {
		args = append(args, *m.Timestamp)
		conditions = append(conditions, "created_at = $"+strconv.Itoa(len(args)))
	}

	args = append(args, userId)
	conditions = append(conditions, "user_id = $"+strconv.Itoa(len(args)))

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += utils.DetermineOrder(m.Order, m.Sort)

	if m.Limit != nil && *m.Limit != "" {
		query += " LIMIT " + *m.Limit
	}

	if m.Offset != nil && *m.Offset != "" {
		query += " OFFSET " + *m.Offset
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (s *TodoStore) SoftUpdate(userId int, id int, model models.TodoUpdateHandlerRequest) (models.Todo, error) {
	oldT, _ := s.getTodo(id, userId)
	var t models.Todo
	err := s.DB.QueryRow("SELECT title, description, done FROM todos WHERE id=$1 AND user_id=$2", id, userId).
		Scan(&t.Title, &t.Description, &t.Done)
	if err != nil {
		return t, err
	}

	if model.Title != nil {
		t.Title = *model.Title
	}
	if model.Description != nil {
		t.Description = *model.Description
	}
	if model.Done != nil {
		t.Done = *model.Done
	}

	err = s.DB.QueryRow(
		"UPDATE todos SET title=$1, description=$2, done=$3 WHERE id=$4 AND user_id=$5 RETURNING id, user_id, title, description, created_at, done",
		t.Title, t.Description, t.Done, id, userId,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)
	if err == nil {
		s.recordHistory(t.ID, userId, oldT, t)
	}
	return t, err
}

func (s *TodoStore) HardUpdate(userId int, id int, model models.TodoUpdateHandlerRequest) (models.Todo, error) {
	oldT, _ := s.getTodo(id, userId)
	var t models.Todo
	err := s.DB.QueryRow(
		"UPDATE todos SET title=$1, description=$2, done=$3 WHERE id=$4 AND user_id=$5 RETURNING id, user_id, title, description, created_at, done",
		model.Title, model.Description, model.Done, id, userId,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)
	if err == nil {
		s.recordHistory(t.ID, userId, oldT, t)
	}
	return t, err
}

func (s *TodoStore) Delete(userId int, id int) (models.Todo, error) {
	oldT, _ := s.getTodo(id, userId)
	var t models.Todo
	err := s.DB.QueryRow(
		"DELETE FROM todos WHERE id=$1 AND user_id=$2 RETURNING id, user_id, title, description, created_at, done",
		id, userId,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)
	if err == nil {
		s.recordHistory(t.ID, userId, oldT, models.Todo{})
	}
	return t, err
}

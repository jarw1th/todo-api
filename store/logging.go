package store

import (
	models "ToDoProject/models"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

func (s *TodoStore) recordHistory(todoID, userID int, oldData, newData models.Todo) error {
	oldB, err := json.Marshal(oldData)
	if err != nil {
		oldB = []byte(fmt.Sprintf("%+v", oldData))
	}

	newB, err := json.Marshal(newData)
	if err != nil {
		newB = []byte(fmt.Sprintf("%+v", newData))
	}

	_, err = s.DB.Exec(
		"INSERT INTO todo_history(todo_id, user_id, old_value, new_value) VALUES($1, $2, $3, $4)",
		todoID, userID, string(oldB), string(newB),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *TodoStore) getTodo(id, userId int) (models.Todo, error) {
	var t models.Todo
	err := s.DB.QueryRow(
		"SELECT id, user_id, title, description, created_at, done FROM todos WHERE user_id=$1 AND id=$2",
		userId, id,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Description, &t.CreatedAt, &t.Done)

	if err != nil {
		return models.Todo{}, err
	}

	return t, nil
}

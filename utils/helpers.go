package utils

import (
	models "ToDoProject/models"
)

func FindFirstById(todos []models.Todo, id int) models.Todo {
	for _, todo := range todos {
		if todo.ID == id {
			return todo
		}
	}
	return models.Todo{}
}

func ChangeFirstById(todos []models.Todo, model models.Todo) []models.Todo {
	for index, todo := range todos {
		if todo.ID == model.ID {
			todos[index] = model
		}
	}
	return todos
}

func DeleteFirstById(todos []models.Todo, id int) ([]models.Todo, models.Todo) {
	var todo models.Todo = models.Todo{}
	for i := 0; i < len(todos); i++ {
		if todos[i].ID == id {
			todo = todos[i]
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	for i := range todos {
		todos[i].ID = i
	}
	return todos, todo
}

func DetermineOrder(order *string, sort *string) string {
	if (order == nil || *order == "") && (sort == nil || *sort == "") {
		return " ORDER BY id"
	}
	var result string = " ORDER BY "
	if sort != nil && *sort != "" {
		result += *sort
	} else {
		result += " created_at"
	}

	if order != nil && *order == "desc" {
		result += " DESC"
	} else {
		result += " ASC"
	}
	return result
}

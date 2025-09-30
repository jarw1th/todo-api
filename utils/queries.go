package utils

import (
	models "ToDoProject/models"
	http "net/http"
	"strconv"
)

func MakeQueriesStruct(r *http.Request) models.TodoQueries {
	var doneBool *bool
	if doneStr := r.URL.Query().Get("done"); doneStr != "" {
		if parsed, err := strconv.ParseBool(doneStr); err == nil {
			doneBool = &parsed
		}
	}
	timestamp := r.URL.Query().Get("created_at")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	model := models.TodoQueries{
		Done:        doneBool,
		Timestamp:   &timestamp,
		Title:       &title,
		Description: &description,
		Order:       &order,
		Sort:        &sort,
		Limit:       &limit,
		Offset:      &offset,
	}
	return model
}

func CheckQueries(m models.TodoQueries) bool {
	return m.Done == nil && *m.Timestamp == "" && *m.Title == "" && *m.Order == "" && *m.Sort == "" && *m.Limit == "" && *m.Offset == "" && *m.Description == ""
}

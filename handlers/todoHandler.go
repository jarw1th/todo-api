package handlers

import (
	"ToDoProject/decode"
	models "ToDoProject/models"
	"ToDoProject/store"
	"ToDoProject/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	mux "github.com/gorilla/mux"
)

type TodoHandler struct {
	Store *store.TodoStore
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.TodoHandlerRequest true "Todo Data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	var req models.TodoHandlerRequest
	if err := decode.DecodeJSONBody(w, r, &req); err != nil {
		if err == decode.ErrEmptyBody {
			decode.JSONError(w, fmt.Errorf("request body cannot be empty"), http.StatusBadRequest)
			return
		}
		decode.JSONError(w, fmt.Errorf("invalid JSON: %w", err), http.StatusBadRequest)
		return
	}

	todo, err := h.Store.Create(userID, req.Title, req.Description)
	if err != nil {
		decode.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// ListTodos godoc
// @Summary List todos
// @Description Get all todos of the authenticated user
// @Tags todos
// @Produce json
// @Param done query bool false "Filter by done"
// @Param title query string false "Filter by title"
// @Param description query string false "Filter by description"
// @Param createdAt query string false "Filter by creation date"
// @Param order query string false "Order by"
// @Param sort query string false "Sort value"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} models.Todo
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	var todos []models.Todo
	var err error

	queries := utils.MakeQueriesStruct(r)

	if utils.CheckQueries(queries) {
		todos, err = h.Store.List(userID)
	} else {
		todos, err = h.Store.FilteredList(userID, queries)
	}

	if err != nil {
		decode.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// PutTodo godoc
// @Summary Update a todo
// @Description Fully update a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.TodoUpdateHandlerRequest true "Todo Data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /todos/{id} [put]
func (h *TodoHandler) PutTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		decode.JSONError(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}

	var req models.TodoUpdateHandlerRequest
	if err := decode.DecodeJSONBody(w, r, &req); err != nil {
		if err == decode.ErrEmptyBody {
			decode.JSONError(w, fmt.Errorf("request body cannot be empty"), http.StatusBadRequest)
			return
		}
		decode.JSONError(w, fmt.Errorf("invalid JSON: %w", err), http.StatusBadRequest)
		return
	}

	if req.Title == nil || req.Description == nil || req.Done == nil {
		decode.JSONError(w, fmt.Errorf("all fields are required for PUT"), http.StatusBadRequest)
		return
	}

	todo, err := h.Store.HardUpdate(userID, id, req)
	if err != nil {
		decode.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// PatchTodo godoc
// @Summary Partially update a todo
// @Description Update only some fields of a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.TodoUpdateHandlerRequest true "Todo Data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /todos/{id} [patch]
func (h *TodoHandler) PatchTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		decode.JSONError(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}

	var req models.TodoUpdateHandlerRequest
	if err := decode.DecodeJSONBody(w, r, &req); err != nil {
		if err == decode.ErrEmptyBody {
			decode.JSONError(w, fmt.Errorf("request body cannot be empty"), http.StatusBadRequest)
			return
		}
		decode.JSONError(w, fmt.Errorf("invalid JSON: %w", err), http.StatusBadRequest)
		return
	}

	todo, err := h.Store.SoftUpdate(userID, id, req)
	if err != nil {
		decode.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, idErr := strconv.Atoi(idStr)
	if idErr != nil {
		decode.JSONError(w, idErr, http.StatusBadRequest)
		return
	}

	todo, err := h.Store.Delete(userID, id)
	if err != nil {
		decode.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

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

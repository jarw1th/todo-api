package handlers

import (
	"ToDoProject/decode"
	token "ToDoProject/jwttoken"
	models "ToDoProject/models"
	"fmt"
	"net/http"
)

func (h *TodoHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decode.DecodeJSONBody(w, r, &req); err != nil {
		status := http.StatusBadRequest
		msg := "invalid JSON"
		if err == decode.ErrEmptyBody {
			msg = "request body cannot be empty"
		}
		decode.JSONError(w, fmt.Errorf(msg+": %w", err), status)
		return
	}

	if req.Username == "" || req.Password == "" {
		decode.JSONError(w, fmt.Errorf("username and password are required"), http.StatusBadRequest)
		return
	}

	userID, err := h.Store.CheckUserCredentials(req.Username, req.Password)
	if err != nil {
		decode.JSONError(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := token.GenerateTokens(userID)
	if err != nil {
		decode.JSONError(w, fmt.Errorf("could not generate tokens: %w", err), http.StatusInternalServerError)
		return
	}

	decode.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *TodoHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := decode.DecodeJSONBody(w, r, &req); err != nil {
		status := http.StatusBadRequest
		msg := "invalid JSON"
		if err == decode.ErrEmptyBody {
			msg = "request body cannot be empty"
		}
		decode.JSONError(w, fmt.Errorf(msg+": %w", err), status)
		return
	}

	if req.Username == "" || req.Password == "" {
		decode.JSONError(w, fmt.Errorf("username and password are required"), http.StatusBadRequest)
		return
	}

	user, err := h.Store.CreateUser(req.Username, req.Password)
	if err != nil {
		decode.JSONError(w, fmt.Errorf("could not create user: %w", err), http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := token.GenerateTokens(user.ID)
	if err != nil {
		decode.JSONError(w, fmt.Errorf("could not generate tokens: %w", err), http.StatusInternalServerError)
		return
	}

	decode.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"user_id":       user.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

package models

import "time"

type Todo struct {
	ID          int       `json:"id"`
	UserId      int       `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Done        bool      `json:"done"`
}

type TodoHandlerRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUpdateHandlerRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type TodoQueries struct {
	Done      *bool
	Timestamp *string
	Title     *string
	Order     *string
	Sort      *string
	Limit     *string
	Offset    *string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int
	Username string
	Password string
}

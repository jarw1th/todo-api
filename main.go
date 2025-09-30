package main

import (
	"ToDoProject/handlers"
	token "ToDoProject/jwttoken"
	recovery "ToDoProject/safety"
	"ToDoProject/store"
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
)

func main() {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=Scooby2011 dbname=todo sslmode=disable"
	todoStore, err := store.NewTodoStore(connStr)
	if err != nil {
		panic(err)
	}
	todoHandler := &handlers.TodoHandler{Store: todoStore}

	r := mux.NewRouter()
	r.Use(recovery.RecoverMiddleware)

	r.HandleFunc("/login", todoHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/register", todoHandler.RegisterHandler).Methods("POST")

	api := r.PathPrefix("/todos").Subrouter()
	api.Use(token.AuthMiddleware)
	api.HandleFunc("", todoHandler.ListTodos).Methods("GET")
	api.HandleFunc("", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/{id}", todoHandler.PutTodo).Methods("PUT")
	api.HandleFunc("/{id}", todoHandler.PatchTodo).Methods("PATCH")
	api.HandleFunc("/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

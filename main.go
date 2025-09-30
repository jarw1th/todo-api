package main

import (
	_ "ToDoProject/docs"
	"ToDoProject/handlers"
	token "ToDoProject/jwttoken"
	recovery "ToDoProject/safety"
	"ToDoProject/store"
	"net/http"

	mux "github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	connStr := recovery.GetDBConnStr()
	todoStore, err := store.NewTodoStore(connStr)
	if err != nil {
		panic(err)
	}
	todoHandler := &handlers.TodoHandler{Store: todoStore}

	r := mux.NewRouter()
	r.Use(recovery.RecoverMiddleware)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/login", todoHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/register", todoHandler.RegisterHandler).Methods("POST")

	api := r.PathPrefix("/todos").Subrouter()
	api.Use(token.AuthMiddleware)
	api.HandleFunc("", todoHandler.ListTodos).Methods("GET")
	api.HandleFunc("", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/{id}", todoHandler.PutTodo).Methods("PUT")
	api.HandleFunc("/{id}", todoHandler.PatchTodo).Methods("PATCH")
	api.HandleFunc("/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	port := recovery.GetPort()
	http.ListenAndServe(":"+port, r)
}

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var todos = []*Todo{{ID: 1, Title: "Todo 1"}, {ID: 2, Title: "Todo 2"}}
var id = 3

func HandleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo := &Todo{ID: id, Title: req.Title}
	id++

	todos = append(todos, todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(*Todo)

	index := slices.Index(todos, todo)
	if index == -1 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos = append(todos[:index], todos[index+1:]...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(*Todo)

	if todo.ID == 3 {
		panic("test")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
}

func HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(*Todo)

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Title = req.Title

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func TodoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		index := slices.IndexFunc(todos, func(todo *Todo) bool { return todo.ID == id })
		if index == -1 {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "todo", todos[index]))

		next.ServeHTTP(w, r)
	})
}

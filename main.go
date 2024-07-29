package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var (
	todos  = []Todo{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/todos", todosHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetTodos(w)
	case "POST":
		handleCreateTodo(w, r)
	case "PUT":
		handleUpdateTodo(w, r)
	case "DELETE":
		handleDeleteTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTodos(w http.ResponseWriter) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(todos)
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	todo.ID = nextID
	nextID++
	todos = append(todos, todo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i := range todos {
		if todos[i].ID == todo.ID {
			todos[i] = todo
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i := range todos {
		if todos[i].ID == todo.ID {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

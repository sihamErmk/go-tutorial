package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Todo struct {
	ID        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

//  global variable to hold todos

// todos is a slice that holds all Todo items in memory.
// It is protected by todoMutex to ensure safe concurrent access.
var (
	todos     = []Todo{}
	todoMutex sync.Mutex
)

func generateRandomId() string {
	return uuid.New().String()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "OK",
		Message: "Api health is OK",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		todoMutex.Lock()
		defer todoMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")

		// Check if the path contains an ID
		path := r.URL.Path
		if len(path) > len("/todo/") {
			// Extract ID from the URL path
			id := path[len("/todo/"):]
			found := false
			for _, t := range todos {
				if t.ID == id {
					json.NewEncoder(w).Encode(t)
					found = true
					break
				}
			}
			if !found {
				http.Error(w, "Todo not found", http.StatusNotFound)
			}
			return
		}

		// Return all todos if no ID is specified
		if len(todos) == 0 {
			http.Error(w, "No todos found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(todos)

	case "POST":
		var todo Todo
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		fmt.Print(body)
		err = json.Unmarshal(body, &todo)
		if err != nil || todo.Task == "" {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		todoMutex.Lock()
		todo.ID = generateRandomId()
		todos = append(todos, todo)
		todoMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	case "PUT":
		var todo Todo
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &todo)
		if err != nil || todo.Task == "" {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		todoMutex.Lock()
		defer todoMutex.Unlock()
		for i, t := range todos {
			if t.ID == todo.ID {
				todos[i] = todo
				json.NewEncoder(w).Encode(todo)
				return
			}
		}
		http.Error(w, "Todo not found", http.StatusNotFound)

	case "DELETE":
		todoMutex.Lock()
		defer todoMutex.Unlock()
		path := r.URL.Path
		if len(path) <= len("/todo/") {
			http.Error(w, "Todo ID is required", http.StatusBadRequest)
			return
		}
		id := path[len("/todo/"):]

		for i, t := range todos {
			if t.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				w.WriteHeader(http.StatusNoContent) // 204 No Content
				return
			}
		}

		http.Error(w, "Todo not found", http.StatusNotFound)

	default:

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	}

}

func main() {
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/todo/", todoHandler)

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("Hello world")
	if err != nil {
		fmt.Println("Error starting server : ", err)
		return
	}

}

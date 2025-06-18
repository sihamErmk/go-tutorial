package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "OK",
		Message: "Api health is OK",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func main() {
	http.HandleFunc("/", healthHandler)
	err := http.ListenAndServe(":8080", nil)
	fmt.Println("Hello world")
	if err != nil {
		fmt.Println("Error starting server : ", err)
		return
	}

}

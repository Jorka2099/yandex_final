package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; harset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)

}

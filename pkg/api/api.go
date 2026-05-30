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
	http.HandleFunc("/api/nextdate", Auth(NextDateHandler))
	http.HandleFunc("/api/task", Auth(taskHandler))
	http.HandleFunc("/api/tasks", Auth(tasksHandler))
	http.HandleFunc("/api/task/done", Auth(doneTaskHandler))
	http.HandleFunc("/api/signin", signinHandler)

}

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "task wasnt found"}, http.StatusBadRequest)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "request body reading error"}, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJSON(w, map[string]string{"error": "JSON Unmarshal error"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Title is empty"}, http.StatusBadRequest)
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]any{}, http.StatusOK)

}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "task was not found"}, http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		writeJSON(w, map[string]any{}, http.StatusOK)
		return
	}

	now := time.Now()

	res, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = db.UpdateDate(res, task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

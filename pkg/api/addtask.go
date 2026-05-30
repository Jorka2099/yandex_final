package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func checkDate(task *db.Task) error {
	var next string
	var err error

	now := time.Now()
	today := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	_, err = time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("wrong date format %w", err)
	}

	isPast := task.Date < today

	// Если есть правило повторения
	if len(task.Repeat) > 0 {
		// Если дата в прошлом - вычисляем следующую
		if isPast {
			next, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("Invalid repeat rule %v", err)
			}
			task.Date = next
		}
		// Если дата сегодня или в будущем - оставляем как есть
		return nil
	}

	// Если нет повторения и дата в прошлом - ставим сегодня
	if isPast {
		task.Date = today
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Request body reading error"}, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
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

	newTaskID, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Failed to add task"}, http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]any{
		"id": strconv.FormatInt(newTaskID, 10),
	}, http.StatusOK)

}

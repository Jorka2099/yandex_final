package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]any{"error": "Failed to load tasks"}, http.StatusBadRequest)
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}



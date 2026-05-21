package db

import (
	"fmt"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask adds task and returnes id
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task = &Task{}
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date =?, title =?, comment=?, repeat=? WHERE id=?"
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task wasnt found")
	}
	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id=?"
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task with this id not found")

	}
	return nil
}

func UpdateDate(next string, task *Task) error {
	task.Date = next

	query := "UPDATE scheduler SET date=? WHERE id=?"
	res, err := DB.Exec(query, task.Date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task with this id not found")
	}
	return nil
}

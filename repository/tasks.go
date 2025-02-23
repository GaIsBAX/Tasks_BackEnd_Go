package repository

import (
	"database/sql"
	model "tasks/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func (r *TaskRepository) GetTasks() ([]model.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) CreateTask(task *model.Task) error {
	return r.DB.QueryRow(
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Description, task.Status,
	).Scan(&task.ID)
}

func (r *TaskRepository) UpdateTask(id int, task *model.Task) error {

	result, err := r.DB.Exec(
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		task.Title, task.Description, task.Status, id,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (r *TaskRepository) DeleteTask(id int) error {
	result, err := r.DB.Exec("DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

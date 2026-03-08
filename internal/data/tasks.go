package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID          int
	CreatedAt   time.Time
	Status      Status
	TaskID      string
	ProjectID   string
	Name        string
	Description string
	Version     int
}

type TaskModel struct {
	DB       *sql.DB
	Projects *ProjectModel
}

func (m *TaskModel) Insert(task Task) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1) get active project and next value
	proj, next, err := m.Projects.GetActiveForUpdate(ctx, tx)
	if err != nil {
		return err
	}

	// 2) build taskId
	taskId := fmt.Sprintf("%s-%05d", proj.ShortName, next)

	// 3) insert task using tx
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO tasks(task_id, project_id, name, status, description)
										values(?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, taskId, proj.ID, task.Name, task.Status.String(), task.Description)
	if err != nil {
		return err
	}

	// increment and update projects next value
	// 4) update project's next_task_value
	if err := m.Projects.IncrementNextTaskValue(ctx, tx, proj.ID, next); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *TaskModel) GetAllTasksForActiveProject() ([]Task, error) {
	var tasks []Task

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `SELECT task_id, name, status, description, version 
										FROM tasks
										WHERE project_id = (SELECT id FROM projects WHERE active = true);`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.TaskID, &t.Name, &t.Status, &t.Description, &t.Version)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m *TaskModel) Update(task Task) error {
	// Increment version but query for previous version to avoid race conditions
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := m.DB.PrepareContext(ctx, `UPDATE tasks SET name = ?, status = ?, description = ?, version = ?
								WHERE id = ?, task_id = ?, version = ?;`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, task.Name, task.Status, task.Description, task.Version+1, task.ID, task.TaskID, task.Version)
	return nil
}

// TODO: DELETE THIS
// HACK: FOR TESTING ONLY
func (m *TaskModel) DeleteAll() {
	_, err := m.DB.Exec(`DELETE FROM tasks;`)
	if err != nil {
		panic(err)
	}
}

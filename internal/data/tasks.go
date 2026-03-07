package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Status int

const (
	DEFINING Status = iota
	TODO
	INPROGRESS
	UNDERREVIEW
	COMPLETED
)

func (s Status) String() string {
	return [...]string{"defining", "todo", "in progress", "under review", "completed"}[s]
}

type Task struct {
	ID          int
	CreatedAt   time.Time
	Status      Status
	TaskID      string
	ProjectID   string
	Name        string
	Description string
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
	// 4 update project's next_task_value
	if err := m.Projects.IncrementNextTaskValue(ctx, tx, proj.ID, next); err != nil {
		return err
	}

	return tx.Commit()
}

// TODO: DELETE THIS
// HACK: FOR TESTING ONLY
func (m *TaskModel) DeleteAll() {
	_, err := m.DB.Exec(`DELETE FROM tasks;`)
	if err != nil {
		panic(err)
	}
}

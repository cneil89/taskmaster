package data

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int
	CreatedAt   time.Time
	TaskID      string
	ProjectID   string
	Name        string
	Description string
}

type TaskModel struct {
	DB *sql.DB
}

// TODO: DELETE THIS
// HACK: FOR TESTING ONLY
func (m *TaskModel) DeleteAll() {
	_, err := m.DB.Exec(`DELETE FROM tasks;`)
	if err != nil {
		panic(err)
	}
}

package data

import (
	"database/sql"
	"time"
)

type Task struct {
	UUID        string
	CreatedAt   time.Time
	TaskID      string
	ProjectID   string
	Name        string
	Description string
}

type TaskModel struct {
	DB *sql.DB
}

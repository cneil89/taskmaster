package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        string
	CreatedAt time.Time
	Name      string
	ShortName string
}

type ProjectModel struct {
	DB *sql.DB
}

func (m *ProjectModel) Insert(name, shortName string) error {

	newID := uuid.New()

	stmt, err := m.DB.Prepare(`INSERT INTO projects (uuid, name, short_name) values (?, ?, ?);`)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = stmt.ExecContext(ctx, newID, name, shortName)
	if err != nil {
		return err
	}

	return nil
}

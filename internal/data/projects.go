package data

import (
	"context"
	"database/sql"
	"time"
)

type Project struct {
	ID        int
	CreatedAt time.Time
	Name      string
	ShortName string
	Active    bool
}

type ProjectModel struct {
	DB *sql.DB
}

// Creates a new project, and sets new project as active
func (m *ProjectModel) Insert(name, shortName string) error {

	stmt, err := m.DB.Prepare(`INSERT INTO projects (name, short_name) values (?, ?) RETURNING id;`)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var projectID int
	err = stmt.QueryRowContext(ctx, name, shortName).Scan(&projectID)
	if err != nil {
		return err
	}

	// Set new project as the active project
	err = m.SetActiveProject(projectID)

	return nil
}

func (m *ProjectModel) GetAllProjects() ([]Project, error) {
	stmt, err := m.DB.Prepare(`SELECT id, name, short_name, active FROM projects;`)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for results.Next() {
		var project Project
		results.Scan(&project.ID, &project.Name, &project.ShortName, &project.Active)
		projects = append(projects, project)
	}

	return projects, nil
}

func (m *ProjectModel) GetActiveProject() (Project, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, `SELECT count(*) FROM projects WHERE active = true;`).Scan(&count)
	if err != nil {
		return Project{}, err
	}
	if count > 1 {
		panic("data integrity issue: more than 1 projects flagged as active.")
	}

	var project Project
	err = m.DB.QueryRowContext(ctx, `SELECT id, name, short_name, active FROM projects WHERE active = true LIMIT 1;`).
		Scan(&project.ID, &project.Name, &project.ShortName, &project.Active)
	if err != nil {
		return Project{}, nil
	}

	return project, nil
}

func (m *ProjectModel) GetActiveForUpdate(ctx context.Context, db DBTX) (Project, int, error) {
	var p Project
	var next int

	row := db.QueryRowContext(ctx, `SELECT id, name, short_name, active, next_task_value FROM projects WHERE active = true;`)
	if err := row.Scan(&p.ID, &p.Name, &p.ShortName, &p.Active, &next); err != nil {
		return Project{}, 0, err
	}

	return p, next, nil
}

// Increment the next task id value
func (m *ProjectModel) IncrementNextTaskValue(ctx context.Context, db DBTX, id int, next int) error {
	// To try an prevent race conditions, use the next value provided to search for the project
	_, err := db.ExecContext(ctx, `UPDATE projects SET next_task_value = ? WHERE id =? AND next_task_value = ?;`, next+1, id, next)

	return err
}

func (m *ProjectModel) SetActiveProject(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE projects SET active = false;`)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `UPDATE projects SET active = true WHERE id = ?`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// TODO: DELETE THIS
// HACK: FOR TESTING ONLY
func (m *ProjectModel) DeleteAll() {
	_, err := m.DB.Exec(`DELETE FROM projects;`)
	if err != nil {
		panic(err)
	}
}

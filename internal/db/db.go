package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func GetDBPath(appName, filename string) (string, error) {
	if appName == "" {
		return "", errors.New("appName Required")
	}

	env := os.Getenv("TASKMASTER_DATA")
	if env != "" {
		env = os.ExpandEnv(env)
		env = strings.TrimSpace(env)
		// expand ~
		if strings.HasPrefix(env, "~") {
			home, err := os.UserHomeDir()
			if err == nil {
				env = filepath.Join(home, strings.TrimPrefix(env, "~"))
			}
		}
		// If env explicitly ends with path separator or exists as a directory => treat as dir
		if strings.HasSuffix(env, string(os.PathSeparator)) {
			p := filepath.Join(env, filename)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return "", err
			}
			return p, nil
		}

		if st, err := os.Stat(env); err == nil && st.IsDir() {
			p := filepath.Join(env, filename)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return "", err
			}
			return p, nil
		}
		// else treat as filepath
		if err := os.MkdirAll(filepath.Dir(env), 0o755); err != nil {
			return "", nil
		}
		return env, nil
	}

	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	case "darwin":
		base = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default:
		base = os.Getenv("XDG_DATA_HOME")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
	}
	dir := filepath.Join(base, appName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(dir, filename), nil
}

func initDB(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("failed to set pragma %q: %w", p, err)
		}
	}

	if err := bootstrapDatabase(db); err != nil {
		return err
	}

	return nil
}

// NOTE: This this be handled by migrations in the future
// TODO: Handle via migrations
func bootstrapDatabase(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY,
			created_at NUMERIC NOT NULL DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			short_name TEXT NOT NULL,
			active BOOLEAN NOT NULL DEFAULT false,
			next_task_num INTEGER NOT NULL DEFAULT 1,
			CONSTRAINT uidx_name_shortname UNIQUE (name)
		);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			created_at NUMERIC NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status TEXT NOT NULL DEFAULT 'defining'
				CHECK(status IN(
					'defining',
					'todo',
					'in-progress',
					'under-review',
					'completed'
				)),
			task_id string NOT NULL UNIQUE,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			CONSTRAINT project_id_fk FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
	}

	for i, st := range stmts {
		_, err := db.Exec(st)
		if err != nil {
			return fmt.Errorf("bootstrap database error: (stmt %d) %s", i, err)
		}
	}

	return nil
}

func OpenDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	var version string
	db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	fmt.Printf("SQLite Version: %s\n", version)

	if err := initDB(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

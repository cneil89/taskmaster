package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cneil89/taskmaster/internal/data"
	"github.com/cneil89/taskmaster/internal/db"
	"github.com/cneil89/taskmaster/internal/vcs"
	"github.com/rivo/tview"
)

type application struct {
	models data.Models
	state  struct {
		selectedRow int

		activeProject     *data.Project
		availableProjects []data.Project

		selectedTask *data.Task
		taskList     []data.Task
		pages        *tview.Pages

		component struct {
			taskTable        *tview.Table
			selectedTaskView *tview.TextView
		}
	}
}

var (
	version = vcs.Version()
	DBPath  = "taskmaster-dev"
)

func main() {
	displayVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Println("Taskmaster CLI")
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	dsn, err := db.GetDBPath(DBPath, "taskm.db")
	if err != nil {
		fmt.Printf("Unable to resolve DB path: %s\n", err)
	}

	db, err := db.OpenDB(dsn)
	if err != nil {
		fmt.Printf("Unable to open database: %s", err.Error())
		os.Exit(1)
	}

	app := application{
		models: data.NewModels(db),
	}

	err = app.Init()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

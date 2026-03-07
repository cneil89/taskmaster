package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cneil89/taskmaster/internal/data"
	"github.com/cneil89/taskmaster/internal/db"
	"github.com/cneil89/taskmaster/internal/vcs"
)

type config struct {
	testing bool
	db      struct {
		dsn string
	}
}

type application struct {
	config config
	models data.Models
}

var version = vcs.Version()

func main() {
	fmt.Println("Taskmaster CLI")
	var cfg config
	flag.BoolVar(&cfg.testing, "testing", false, "toggle testing WARNING: WILL DELETE ALL DATA")

	displayVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	dsn, err := db.GetDBPath("taskmaster", "taskm.db")
	if err != nil {
		fmt.Printf("Unable to resolve DB path: %s\n", err)
	}

	cfg.db.dsn = dsn

	db, err := db.OpenDB(cfg.db.dsn)
	if err != nil {
		fmt.Printf("Unable to open database: %s", err.Error())
		os.Exit(1)
	}

	app := application{
		config: cfg,
		models: data.NewModels(db),
	}

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

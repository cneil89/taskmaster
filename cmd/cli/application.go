package main

import "fmt"

func (app *application) Run() error {
	fmt.Println("Running Insert Statement: Taskmaster, taskm")
	err := app.models.Projects.Insert("Taskmaster", "taskm")
	if err != nil {
		return err
	}

	return nil
}

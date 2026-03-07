package main

import (
	"fmt"

	"github.com/cneil89/taskmaster/internal/data"
)

func (app *application) Run() error {
	if app.config.testing {
		// HACK: Delete all data in db so it can be repopulated
		fmt.Println("\n- Deleting all data")
		app.models.Tasks.DeleteAll()
		app.models.Projects.DeleteAll()
	}
	app.insertProject("Taskmaster", "taskm")
	app.addTasks(5)
	app.insertProject("Testing 1", "test1")
	app.addTasks(4)
	app.insertProject("Testing 2", "test2")
	app.addTasks(8)

	app.printAllProjects()
	app.getActiveProject()
	app.setActiveProject(1) // Project: Taskmaster

	app.printAllProjects()
	app.getActiveProject()

	app.addTasks(12)

	return nil
}

// TODO: DELETE These
// NOTE: Helper functions for testing
func (app *application) addTasks(n int) {
	fmt.Printf("- Running Insert Tasks (count: %d)\n", n)
	for i := range n {
		task := data.Task{
			Name:        fmt.Sprintf("test %d", i),
			Description: "Test task",
			Status:      data.Status(i % 5),
		}

		if err := app.models.Tasks.Insert(task); err != nil {
			panic(err)
		}
	}
}
func (app *application) insertProject(name, sname string) {
	fmt.Printf("- Running Insert Project: %s, %s\n", name, sname)
	if err := app.models.Projects.Insert(name, sname); err != nil {
		panic(err)
	}
}

func (app *application) printAllProjects() {
	fmt.Println("- Getting All Projects")
	projects, err := app.models.Projects.GetAllProjects()
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		active := " "
		if project.Active {
			active = "*"
		}
		fmt.Printf("\t- Project: %s %d %s, %s\n", active, project.ID, project.Name, project.ShortName)
	}
}

func (app *application) setActiveProject(id int) {
	fmt.Printf("- Setting active project: id = %d\n", id)
	err := app.models.Projects.SetActiveProject(id)
	if err != nil {
		panic(err)
	}
}

func (app *application) getActiveProject() {
	fmt.Println("Getting active project")
	prj, err := app.models.Projects.GetActiveProject()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\t- Active Project: %d %q %q\n", prj.ID, prj.Name, prj.ShortName)
}

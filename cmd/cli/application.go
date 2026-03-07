package main

import "fmt"

func (app *application) Run() error {
	if app.config.testing {
		// HACK: Delete all data in db so it can be repopulated
		fmt.Println("\nDELETING DATA")
		app.models.Tasks.DeleteAll()
		app.models.Projects.DeleteAll()
	}
	app.insertProject("Taskmaster", "taskm")
	app.insertProject("Testing 1", "test1")
	app.insertProject("Testing 2", "test2")

	app.printAllProjects()
	app.getActiveProject()

	app.setActiveProject(1) // Project: Taskmaster

	app.printAllProjects()
	app.getActiveProject()

	return nil
}

// TODO: DELETE These
// INFO: Helper functions for testing
func (app *application) insertProject(name, sname string) {
	fmt.Printf("Running Insert Statement: %s, %s\n", name, sname)
	err := app.models.Projects.Insert(name, sname)
	if err != nil {
		panic(err)
	}
}

func (app *application) printAllProjects() {
	fmt.Println("\nGetting All Projects")
	projects, err := app.models.Projects.GetAllProjects()
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		active := " "
		if project.Active {
			active = "*"
		}
		fmt.Printf("Project: %s %d %s, %s\n", active, project.ID, project.Name, project.ShortName)
	}
}

func (app *application) setActiveProject(id int) {
	fmt.Printf("\nSetting active project: id = %d", id)
	err := app.models.Projects.SetActiveProject(id)
	if err != nil {
		panic(err)
	}
}

func (app *application) getActiveProject() {
	fmt.Println("\nGetting active project")
	prj, err := app.models.Projects.GetActiveProject()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Active Project: %d %q %q\n", prj.ID, prj.Name, prj.ShortName)
}

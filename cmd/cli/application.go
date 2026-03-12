package main

import (
	"fmt"

	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *application) Init() error {

	var err error

	app.state.activeProject, err = app.models.Projects.GetActiveProject()
	if err != nil {
		return err
	}

	app.state.availableProjects, err = app.models.Projects.GetAllProjects()
	if err != nil {
		return err
	}

	app.state.taskList, err = app.models.Tasks.GetAllTasksForActiveProject()
	if err != nil {
		return err
	}

	if len(app.state.taskList) == 0 {
		app.state.selectedTask = nil
	} else {
		app.state.selectedTask = &app.state.taskList[app.state.selectedRow]
	}

	app.state.selectedTaskView = tview.NewTextView().SetText("")
	app.state.selectedTaskView.SetDrawFunc(
		func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			var text string
			str := "\n%-14s %s\n%-14s %s\n%-14s %s\n%-14s %s\n%-14s %s"

			activeProj := &data.Project{Name: "", ShortName: ""}
			if app.state.activeProject != nil {
				activeProj = app.state.activeProject
			}

			if app.state.selectedTask == nil {
				text = fmt.Sprintf(
					str,
					"Project:", activeProj.Name,
					"Task ID:", "",
					"Name:", "",
					"Status:", "",
					"Description:", "",
				)
			} else {
				text = fmt.Sprintf(
					str,
					"Project:", activeProj.Name,
					"Task ID:", app.state.selectedTask.TaskID,
					"Name:", app.state.selectedTask.Name,
					"Status:", app.state.selectedTask.Status,
					"Description:", app.state.selectedTask.Description,
				)
			}

			app.state.selectedTaskView.SetText(text)

			return x, y, width, height
		})

	return nil
}

func (app *application) Run() error {
	if app.config.testing {
		// HACK: Delete all data in db so it can be repopulated
		fmt.Println("\n- Deleting all data")
		app.models.Tasks.DeleteAll()
		app.models.Projects.DeleteAll()
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

		app.printActiveProjectTasks()
		app.setActiveProject(3) // Project: Test2
		app.printActiveProjectTasks()
		app.updateTask(2)
		app.printActiveProjectTasks()
	}

	app.buildTaskTable()
	app.state.pages = tview.NewPages()

	textView := tview.NewTextView().
		SetText("\n" + LOGO).SetTextColor(tcell.ColorDarkCyan).SetTextAlign(tview.AlignCenter)

	legendView := tview.NewTextView().
		SetText("p: Select Project | P: New Project  |  t: Add Task  |  +/-: Quick Status Update").
		SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(textView, 0, 2, false).
		AddItem(app.state.selectedTaskView, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(flex, 0, 1, false).
		AddItem(app.state.pages, 0, 3, true).
		AddItem(legendView, 1, 0, false)

	app.state.pages.AddPage("taskList", app.state.component.taskTable, true, true)
	if app.state.activeProject == nil {
		app.showCreateProjectModal()
	}

	frontend := tview.NewApplication().SetRoot(layout, true)
	return frontend.Run()
}

// TODO: DELETE These
// NOTE: Helper functions for testing

func (app *application) updateTask(id int) {
	fmt.Println("- Updating Task")
	tasks, err := app.models.Tasks.GetAllTasksForActiveProject()
	if err != nil {
		panic(err)
	}

	activeTask := tasks[id]
	activeTask.Description = "Batman was here"

	err = app.models.Tasks.Update(activeTask)
	if err != nil {
		panic(err)
	}

}

func (app *application) printActiveProjectTasks() {
	fmt.Println("- Printing Tasks for Active Project")
	var err error
	app.state.taskList, err = app.models.Tasks.GetAllTasksForActiveProject()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\t   %10s | %15s | %14s | %s | %s\n", "TaskID", "TaskName", "TaskStatus", "v", "TaskDescription")
	fmt.Printf("\t  -%10s-|-%15s-|-%14s-|-%s-|-%s\n", "----------", "---------------", "--------------", "-", "--------")
	for _, task := range app.state.taskList {
		fmt.Printf("\t- %10s | %15s | %14s | %d | %s\n", task.TaskID, task.Name, task.Status, task.Version, task.Description)
	}

}

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
	var err error
	app.state.availableProjects, err = app.models.Projects.GetAllProjects()
	if err != nil {
		panic(err)
	}

	for _, project := range app.state.availableProjects {
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
	fmt.Println("- Getting active project")
	var err error
	app.state.activeProject, err = app.models.Projects.GetActiveProject()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\t- Active Project: %d %q %q\n", app.state.activeProject.ID,
		app.state.activeProject.Name, app.state.activeProject.ShortName)
}

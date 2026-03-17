package main

import (
	"fmt"

	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *application) Init() error {

	app.state.pages = tview.NewPages()

	err := app.updateState()
	if err != nil {
		return err
	}

	app.state.component.selectedTaskView = tview.NewTextView().SetText("")
	app.state.component.selectedTaskView.SetDrawFunc(
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

			app.state.component.selectedTaskView.SetText(text)

			return x, y, width, height
		})

	return nil
}

func (app *application) updateState() error {

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

	if len(app.state.taskList) <= app.state.selectedRow {
		app.state.selectedTask = nil
		app.state.selectedRow = 0

		if len(app.state.taskList) > 0 {
			app.state.selectedTask = &app.state.taskList[app.state.selectedRow]
		}
	} else {
		app.state.selectedTask = &app.state.taskList[app.state.selectedRow]
	}

	app.buildTaskTable()
	return nil
}

func (app *application) Run() error {

	logoView := tview.NewTextView().
		SetText(LOGO).SetTextColor(tcell.ColorDarkCyan).SetTextAlign(tview.AlignCenter)

	legendView := tview.NewTextView().
		SetText("p: Select Project | P: New Project | t: Add Task | Enter: Edit Task | +/-: Quick Status Update | ?: Info").
		SetTextAlign(tview.AlignCenter)

	rowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(logoView, 0, 2, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(rowFlex, 0, 2, false).
		AddItem(app.state.component.selectedTaskView, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(flex, 0, 1, false).
		AddItem(app.state.pages, 0, 2, true).
		AddItem(legendView, 1, 0, false)

	err := app.updateState()
	if err != nil {
		app.notifyError(err)
	}

	if app.state.activeProject == nil {
		app.showCreateProjectModal()
	}

	frontend := tview.NewApplication().SetRoot(layout, true)
	return frontend.Run()
}

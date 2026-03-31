package main

import (
	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *application) buildTaskTable() {

	tableHeaders := []string{"Task Id", "Task Name", "Status", "Description"}

	app.state.component.taskTable = tview.NewTable().
		SetFixed(1, 0).
		SetSeparator(tview.Borders.Vertical).
		SetSelectedStyle(
			tcell.StyleDefault.
				Background(tcell.ColorDarkCyan).
				Foreground(tcell.ColorWhite),
		)

	for column, header := range tableHeaders {
		expand := 1
		width := 0
		if column == 3 {
			expand = 0
			width = DESCRIPTION_TRUNCATE
		}
		app.state.component.taskTable.SetCell(0, column,
			&tview.TableCell{
				Text:          header,
				Color:         tcell.ColorDarkCyan,
				Align:         tview.AlignCenter,
				NotSelectable: true,
				Expansion:     expand,
				MaxWidth:      width,
			})
	}

	for row, task := range app.state.taskList {
		for col := range len(tableHeaders) {
			align := tview.AlignLeft
			expand := 1
			width := 0
			text := taskCellValue(task, col)
			if col == 2 {
				align = tview.AlignCenter
			}
			if col == 3 {
				expand = 0
				width = DESCRIPTION_TRUNCATE
				text = truncate(taskCellValue(task, col), DESCRIPTION_TRUNCATE)
			}
			app.state.component.taskTable.SetCell(row+1, col,
				&tview.TableCell{
					Text:      text,
					Color:     tcell.ColorWhite,
					Align:     align,
					Expansion: expand,
					MaxWidth:  width,
				},
			)
		}
	}

	app.state.component.taskTable.SetSelectionChangedFunc(func(row, col int) {
		if row <= 0 {
			app.state.selectedRow = -1
			app.state.selectedTask = nil
			return
		}

		idx := row - 1
		app.state.selectedRow = idx
		app.state.selectedTask = &app.state.taskList[idx]
	})

	app.state.component.taskTable.SetSelectedFunc(func(row, column int) {
		app.editTask()
	})

	if len(app.state.taskList) > 0 {
		app.state.component.taskTable.SetSelectable(true, false)

		if app.state.selectedRow < 0 {
			app.state.selectedRow = 0
		}

		app.state.component.taskTable.Select(app.state.selectedRow+1, 0)
	}

	app.state.component.taskTable.
		SetBorder(true).
		SetTitle("Tasks").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'P':
				app.showCreateProjectModal()
				return nil
			case 'p':
				app.selectDifferentProject()
				return event
			case 't':
				app.createNewTask()
				return event
			case '+', '=':
				app.incrementTaskStatus(*app.state.selectedTask)
				return event
			case '-', '_':
				app.decrementTaskStatus(*app.state.selectedTask)
				return event
			}

			return event
		})
	if app.state.pages.GetPage("taskList") != nil {
		app.state.pages.RemovePage("taskList")
	}
	app.state.pages.AddPage("taskList", app.state.component.taskTable, true, true)
}

func (app *application) editTask() {
	tmp := data.Task{
		ID:          app.state.selectedTask.ID,
		Version:     app.state.selectedTask.Version,
		Name:        app.state.selectedTask.Name,
		TaskID:      app.state.selectedTask.TaskID,
		Status:      app.state.selectedTask.Status,
		Description: app.state.selectedTask.Description,
	}
	form := tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkCyan).
		SetButtonBackgroundColor(tcell.ColorSlateGrey).
		AddInputField("Name:", app.state.selectedTask.Name, 0, nil, func(v string) {
			tmp.Name = v
		}).
		AddDropDown("Status", []string{
			data.DEFINING.String(),
			data.READY.String(),
			data.INPROGRESS.String(),
			data.UNDERREVIEW.String(),
			data.COMPLETED.String(),
		},
			int(app.state.selectedTask.Status),
			func(option string, index int) {
				val, err := data.ParseStatus(option)
				if err != nil {
					app.notifyError(err)
				}
				tmp.Status = val
			}).
		AddTextArea("Description", app.state.selectedTask.Description, 0, 0, 0, func(changed string) {
			tmp.Description = changed
		}).
		AddButton("Save", func() {
			err := app.models.Tasks.Update(tmp)
			if err != nil {
				app.notifyError(err)
			}
			app.state.selectedTask = &tmp

			err = app.updateState()
			if err != nil {
				app.notifyError(err)
			}
			app.state.pages.RemovePage("modal")
		}).
		AddButton("Cancel", func() {
			app.state.pages.RemovePage("modal")
		})

	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)

	showModal(app, 60, 16, form)
}

func (app *application) createNewTask() {
	tmp := data.Task{
		ProjectID: app.state.activeProject.ID,
	}

	form := tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkCyan).
		SetButtonBackgroundColor(tcell.ColorSlateGrey).
		AddInputField("Task Name:", "", 0, nil, func(v string) {
			tmp.Name = v
		}).
		AddDropDown("Status", []string{
			data.DEFINING.String(),
			data.READY.String(),
			data.INPROGRESS.String(),
			data.UNDERREVIEW.String(),
			data.COMPLETED.String(),
		},
			int(tmp.Status),
			func(option string, index int) {
				val, err := data.ParseStatus(option)
				if err != nil {
					app.notifyError(err)
				}
				tmp.Status = val
			}).
		AddTextArea("Description", "", 0, 0, 0, func(changed string) {
			tmp.Description = changed
		}).
		AddButton("Save", func() {
			err := app.models.Tasks.Insert(tmp)
			if err != nil {
				app.notifyError(err)
				return
			}

			err = app.updateState()
			if err != nil {
				app.notifyError(err)
				return
			}

			app.state.pages.RemovePage("modal")
		}).
		AddButton("Cancel", func() {
			app.state.pages.RemovePage("modal")
		})

	form.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle("Create New Project").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'p':
				return event
			}

			return event
		})

	showModal(app, 60, 20, form)
}

func (app *application) incrementTaskStatus(task data.Task) {
	if app.state.selectedTask == nil {
		return
	}

	tmp := task
	if tmp.Status < data.COMPLETED {
		tmp.Status++
	}

	err := app.models.Tasks.Update(tmp)
	if err != nil {
		app.notifyError(err)
		return
	}

	app.updateState()

}

func (app *application) decrementTaskStatus(task data.Task) {
	if app.state.selectedTask == nil {
		return
	}

	tmp := task
	if tmp.Status > data.DEFINING {
		tmp.Status--
	}

	err := app.models.Tasks.Update(tmp)
	if err != nil {
		app.notifyError(err)
		return
	}

	app.updateState()
}

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
			tcell.StyleDefault.Background(tcell.ColorDarkCyan).Foreground(tcell.ColorWhite),
		)
	for column, header := range tableHeaders {
		expand := 1
		if column == 3 {
			expand = 3
		}
		app.state.component.taskTable.SetCell(0, column,
			&tview.TableCell{
				Text:          header,
				Color:         tcell.ColorDarkCyan,
				Align:         tview.AlignCenter,
				NotSelectable: true,
				Expansion:     expand,
			})
	}

	for row, task := range app.state.taskList {
		for col := range len(tableHeaders) {
			align := tview.AlignLeft
			expand := 1
			if col == 2 {
				align = tview.AlignCenter
			}
			if col == 3 {
				expand = 3
			}
			app.state.component.taskTable.SetCell(row+1, col,
				&tview.TableCell{
					Text:      taskCellValue(task, col),
					Color:     tcell.ColorWhite,
					Align:     align,
					Expansion: expand,
				},
			)
		}
	}

	app.state.component.taskTable.SetSelectionChangedFunc(func(row, col int) {
		if row == 0 {
			return
		}

		app.state.selectedRow = row
		app.state.selectedTask = &app.state.taskList[row-1]
	})

	app.state.component.taskTable.SetSelectedFunc(func(row, column int) {
		app.editTask()
	})

	if len(app.state.taskList) > 0 {
		app.state.component.taskTable.
			SetSelectable(true, false).
			Select(app.state.selectedRow, 0)

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
			}

			return event
		})

}

func (app *application) editTask() {
	tmp := data.Task{
		ID:          app.state.selectedTask.ID,
		Version:     app.state.selectedTask.Version,
		Name:        app.state.activeProject.Name,
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
			data.TODO.String(),
			data.INPROGRESS.String(),
			data.UNDERREVIEW.String(),
			data.COMPLETED.String(),
		},
			int(app.state.selectedTask.Status),
			func(option string, index int) {
				val, err := data.ParseStatus(option)
				if err != nil {
					panic(err)
				}
				tmp.Status = val
			}).
		AddTextArea("Description", app.state.selectedTask.Description, 0, 0, 0, func(changed string) {
			tmp.Description = changed
		}).
		AddButton("Save", func() {
			err := app.models.Tasks.Update(tmp)
			if err != nil {
				// TODO: Need to do something that signifies that the update failed
				panic(err)
			}
			app.state.selectedTask = &tmp

			app.state.taskList, err = app.models.Tasks.GetAllTasksForActiveProject()
			if err != nil {
				panic(err)
			}

			app.buildTaskTable()
			app.state.pages.RemovePage("modal")
			app.state.pages.RemovePage("taskList")
			app.state.pages.AddPage("taskList", app.state.component.taskTable, true, true)

		}).
		AddButton("Cancel", func() {
			app.state.pages.RemovePage("modal")
		})

	form.SetBorder(true).SetTitleAlign(tview.AlignCenter)

	showModal(app, 60, 16, form)
}

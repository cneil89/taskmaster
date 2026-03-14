package main

import (
	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showCreateProjectModal presents a modal containing a form to input information to create
// a new project
func (app *application) showCreateProjectModal() {
	tmp := data.Project{}
	form := tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorDarkCyan).
		SetButtonBackgroundColor(tcell.ColorSlateGrey).
		AddInputField("Project Name:", "", 0, nil, func(v string) {
			tmp.Name = v
		}).
		AddInputField("Short Name:", "", 0, nil, func(v string) {
			tmp.ShortName = v
		}).
		AddButton("Save", func() {
			// err := app.models.Tasks.Update(tmp)
			err := app.models.Projects.Insert(tmp.Name, tmp.ShortName)
			if err != nil {
				// TODO: Gracefully fail, and inform user
				panic(err)
			}

			err = app.updateState()
			if err != nil {
				// TODO: Gracefully fail, and inform user
				panic(err)
			}

			app.state.pages.RemovePage("modal")
			app.state.pages.RemovePage("taskList")
			app.state.pages.AddPage("taskList", app.state.component.taskTable, true, true)
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

	showModal(app, 60, 10, form)
}

// function to select a different project.
func (app *application) selectDifferentProject() {
	var err error
	// Get all projects
	app.state.availableProjects, err = app.models.Projects.GetAllProjects()
	if err != nil {
		// TODO: Gracefully fail, and inform user
		panic(err)
	}

	// make table
	tableHeaders := []string{"Active", "Project"}

	projectTable := tview.NewTable().
		SetFixed(1, 0).
		SetSelectedStyle(
			tcell.StyleDefault.Background(tcell.ColorDarkCyan).Foreground(tcell.ColorWhite),
		)
	for column, header := range tableHeaders {
		projectTable.SetCell(0, column,
			&tview.TableCell{
				Text:          header,
				Color:         tcell.ColorDarkCyan,
				Align:         tview.AlignCenter,
				NotSelectable: true,
				Expansion:     1,
			},
		)
	}

	for row, project := range app.state.availableProjects {
		activeMarker := ""
		color := tcell.ColorWhite
		if project.Active {
			activeMarker = "*"
			color = tcell.ColorDarkGoldenrod
			projectTable.Select(row+1, 0)
		}
		projectTable.SetCell(row+1, 0,
			&tview.TableCell{
				Text:      activeMarker,
				Color:     color,
				Align:     tview.AlignCenter,
				Expansion: 1,
			},
		)

		projectTable.SetCell(row+1, 1,
			&tview.TableCell{
				Text:      project.Name,
				Color:     color,
				Align:     tview.AlignLeft,
				Expansion: 1,
			},
		)
	}

	projectTable.SetSelectedFunc(func(row, column int) {
		err := app.models.Projects.SetActiveProject(app.state.availableProjects[row-1].ID)
		if err != nil {
			// TODO: Gracefully fail, and inform user
			panic(err)
		}

		err = app.updateState()
		if err != nil {
			// TODO: Gracefully fail, and inform user
			panic(err)
		}

		app.state.pages.RemovePage("modal")
		app.state.pages.RemovePage("taskList")
		app.state.pages.AddPage("taskList", app.state.component.taskTable, true, true)
	})

	if len(app.state.availableProjects) > 0 {
		projectTable.
			SetSelectable(true, false)
	}

	projectTable.
		SetBorder(true).
		SetTitle("Projects")

	// display in modal
	modalHeight := min(len(app.state.availableProjects)+3, 15)
	showModal(app, 60, modalHeight, projectTable)

}

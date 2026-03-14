package main

import (
	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const LOGO = `████████╗ █████╗ ███████╗██╗  ██╗███╗   ███╗ █████╗ ███████╗████████╗███████╗██████╗
 ╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝████╗ ████║██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗
    ██║   ███████║███████╗█████╔╝ ██╔████╔██║███████║███████╗   ██║   █████╗  ██████╔╝
    ██║   ██╔══██║╚════██║██╔═██╗ ██║╚██╔╝██║██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗
    ██║   ██║  ██║███████║██║  ██╗██║ ╚═╝ ██║██║  ██║███████║   ██║   ███████╗██║  ██║
    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝`

func taskCellValue(t data.Task, col int) string {
	switch col {
	case 0:
		return t.TaskID
	case 1:
		return t.Name
	case 2:
		return t.Status.String()
	case 3:
		return t.Description
	default:
		return ""
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-3] + "..."
}

func centered(app *application, w, h int, p tview.Primitive) tview.Primitive {
	modalFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, h, 1, true).
				AddItem(nil, 0, 1, false),
			w,
			1,
			true,
		).
		AddItem(nil, 0, 1, false)

	modalFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.state.pages.RemovePage("modal")
			return nil
		}

		return event
	})

	return modalFlex
}

func showModal(app *application, w, h int, p tview.Primitive) {
	modal := centered(app, w, h, p)

	app.state.pages.AddPage("modal", modal, true, true)
}

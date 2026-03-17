package main

import (
	"strings"
	"unicode/utf8"

	"github.com/cneil89/taskmaster/internal/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DESCRIPTION_TRUNCATE = 85

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

func (app *application) notifyError(err error) {
	width := 75
	text := err.Error() + "\n\n" + centeredString("Press ESC to close window", width)
	wrappedLines := computeWrappedLines(text, width-4)
	padding := 3
	desiredHeight := wrappedLines + padding
	maxHeight := 20

	height := desiredHeight
	overflow := false
	if height > maxHeight {
		height = maxHeight
		overflow = true
	}

	errTextView := tview.NewTextView()
	errTextView.SetTextColor(tcell.ColorRed)
	errTextView.SetText(text)
	errTextView.SetWrap(true)
	errTextView.SetBorder(true)
	errTextView.SetTitle("ERROR")
	errTextView.SetBorderColor(tcell.ColorRed)

	if overflow {
		errTextView.SetScrollable(true)
	}

	showModal(app, width, height, errTextView)
}

func centeredString(s string, width int) string {
	if len(s) > width {
		return s
	}
	var paddedStr string
	padding := (width - len(s)) / 2

	paddedStr = strings.Repeat(" ", padding) + s
	return paddedStr
}

func computeWrappedLines(s string, width int) int {
	if width <= 0 {
		return 0
	}
	lines := 0

	for rawLine := range strings.SplitSeq(s, "\n") {
		rcount := utf8.RuneCountInString(rawLine)
		if rcount == 0 {
			lines++
			continue
		}
		lines += (rcount + width - 1) / width
	}

	return lines
}

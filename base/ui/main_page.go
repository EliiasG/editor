package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/state"
)

func NewMainPage(app *state.EditorApp) fyne.CanvasObject {
	exp := NewProjectExplorer(app)
	split := container.NewHSplit(exp, widget.NewLabel("WIP2"))
	split.SetOffset(0.25)
	return split
}

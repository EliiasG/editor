package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/eliiasg/editor/base/state"
)

func NewMainPage(app *state.EditorApp) fyne.CanvasObject {
	exp := NewProjectExplorer(app)
	tabs := NewTabArea(app)
	split := container.NewHSplit(exp, tabs)
	split.SetOffset(0.25)
	return split
}

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/eliiasg/editor/base/state"
)

const DoubleClickTimeMS int64 = 500

func NewRoot(app *state.EditorApp) fyne.CanvasObject {
	root := container.NewMax()
	root.Objects = []fyne.CanvasObject{
		NewStartPage(app, func(i int) {
			app.ProjectManager().OpenProject(i)
			root.Objects = []fyne.CanvasObject{NewMainPage(app)}
		}),
	}
	return root
}

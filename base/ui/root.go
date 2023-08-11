package ui

import (
	"fyne.io/fyne/v2"
	"github.com/eliiasg/editor/base/state"
)

func GetRoot(app *state.EditorApp) fyne.CanvasObject {

	return GetStartPage(app, func(i int) {
		println(app.GetProjectManager().GetRecentProjects()[i])
	})
}

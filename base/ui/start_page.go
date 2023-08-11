package ui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/state"
)

func GetStartPage(app *state.EditorApp, open func(int)) fyne.CanvasObject {
	projectButtonBox := container.NewVBox()
	removeButtonBox := container.NewVBox()
	return container.NewCenter(container.NewHBox(projectButtonBox, removeButtonBox))
}

func resetProjectList(app *state.EditorApp, open func(int), projectButtonObjects, removeButtonObjects *fyne.Container) {
	//projectButtonObjects = GenerateProjectButtons(app, open)
}

func GenerateProjectButtons(app *state.EditorApp, open func(int)) []fyne.CanvasObject {
	recent := []string{"hey", "world"} //app.GetProjectManager().GetRecentProjects()
	elements := make([]fyne.CanvasObject, 0, len(recent))
	for _, path := range recent {
		button := widget.NewButton(
			path,
			func() {
				tryOpen(path, app.GetProjectManager(), app.GetMainWindow(), open)
			},
		)
		elements = append(elements, button)
	}
	return elements
}

func GenerateRemoveButtons(app *state.EditorApp, remove func(int)) []fyne.CanvasObject {
	recent := []string{"hey", "world"} //app.GetProjectManager().GetRecentProjects()
	elements := make([]fyne.CanvasObject, 0, len(recent))
	for i, path := range recent {
		button := widget.NewButton(path, func() { remove(i) })
		elements = append(elements, container.NewHBox(button))
	}
	return elements
}

func tryOpen(path string, manager *state.ProjectManager, win fyne.Window, open func(int)) {
	//new := !slices.Contains(recent, path)
	conf := dialog.NewConfirm(
		"Open project",
		"Do you wnat to open \""+path+"\"?",
		func(b bool) {
			if b {
				idx := slices.Index(manager.GetRecentProjects(), path)
				if idx == -1 {
					manager.AddProject(path)
					idx = 0
				}
				open(idx)
			}
		},
		win,
	)
	conf.Show()
}

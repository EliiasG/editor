package ui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/state"
)

func GetStartPage(app *state.EditorApp, open func(int)) fyne.CanvasObject {
	projectButtonBox := container.NewVBox()
	removeButtonBox := container.NewVBox()
	rst := func() {
		resetProjectList(app, open, projectButtonBox, removeButtonBox)
	}
	rst()
	lst := container.NewVScroll(
		container.NewBorder(
			nil,
			nil,
			nil,
			removeButtonBox,
			projectButtonBox,
		),
	)
	return container.NewBorder(
		container.NewHBox(
			getAddButton(app, rst),
		),
		nil,
		nil,
		nil,
		container.NewPadded(lst),
	)
}

func resetProjectList(app *state.EditorApp, open func(int), projectButtonContainer, removeButtonContainer *fyne.Container) {
	projectButtonContainer.Objects = generateProjectButtons(app, open)
	removeButtonContainer.Objects = generateRemoveButtons(app, func(idx int) {
		projman := app.GetProjectManager()
		projman.RemoveProject(idx)
		projman.Save(app)
		resetProjectList(app, open, projectButtonContainer, removeButtonContainer)
	})
}

func getAddButton(app *state.EditorApp, reset func()) fyne.CanvasObject {
	return widget.NewButton("Add Project", func() {
		openDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			if lu == nil {
				return
			}
			projman := app.GetProjectManager()
			projman.AddProject(lu.Path())
			projman.Save(app)
			reset()
		}, app.GetMainWindow())
		openDialog.Show()
	})
}

func generateProjectButtons(app *state.EditorApp, open func(int)) []fyne.CanvasObject {
	recent := app.GetProjectManager().GetRecentProjects()
	elements := make([]fyne.CanvasObject, 0, len(recent))
	for _, path := range recent {
		pathCopy := path
		button := widget.NewButton(
			shortenPath(pathCopy),
			func() {
				tryOpen(pathCopy, app.GetProjectManager(), app.GetMainWindow(), open)
			},
		)
		button.Alignment = widget.ButtonAlignLeading
		elements = append(elements, container.NewHBox(button))
	}
	return elements
}

func shortenPath(path string) string {
	if len(path) >= 80 {
		return path[:79] + "..."
	} else {
		return path
	}
}

func generateRemoveButtons(app *state.EditorApp, remove func(idx int)) []fyne.CanvasObject {
	recent := app.GetProjectManager().GetRecentProjects()
	elements := make([]fyne.CanvasObject, 0, len(recent))
	for i, path := range recent {
		//i is updated, using i will always remove last project
		pathCopy := path
		j := i
		button := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
			conf := dialog.NewConfirm("Remove project", "Do you want to remove \""+shortenPath(pathCopy)+"\" from the project list?", func(b bool) {
				if b {
					remove(j)
				}
			}, app.GetMainWindow())
			conf.Show()
		})
		elements = append(elements, container.NewHBox(button))
	}
	return elements
}

func tryOpen(path string, manager *state.ProjectManager, win fyne.Window, open func(int)) {
	//new := !slices.Contains(recent, path)
	conf := dialog.NewConfirm(
		"Open project",
		"Do you wnat to open \""+shortenPath(path)+"\"?",
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

package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/state"
)

type ClickableLabel struct {
	widget.Label
	Clicked      func()
	RightClicked func()
}

func (c *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	if c.Clicked != nil {
		c.Clicked()
	}
}

func (c *ClickableLabel) TappedSecondary(_ *fyne.PointEvent) {
	if c.RightClicked != nil {
		c.RightClicked()
	}
}

func NewClickableLabel(text string, click func()) fyne.CanvasObject {
	cl := &ClickableLabel{Label: widget.Label{Text: text}}
	cl.Clicked = click
	cl.ExtendBaseWidget(cl)
	return cl
}

func NewMainPage(app *state.EditorApp) fyne.CanvasObject {
	var cl fyne.CanvasObject
	var rect fyne.CanvasObject
	var reloadExp func()
	exp, reloadExp := NewProjectExplorer(app, func(path string) {
		cl.Hide()
		rect.Hide()
		if app.ProjectManager().IsRequesting() {
			app.ProjectManager().ConfirmFileSelection(path)
			reloadExp()
		}
	})
	tabs := NewTabArea(app)

	rect = canvas.NewRectangle(color.RGBA{0, 0, 0, 150})
	cl = NewClickableLabel("", func() {
		cl.Hide()
		rect.Hide()
		app.ProjectManager().CancelFileSelection()
		reloadExp()
	})
	app.ProjectManager().UIRequestFile = func() {
		reloadExp()
		cl.Show()
		rect.Show()
	}

	split := container.NewHSplit(container.NewHScroll(exp), container.NewMax(
		tabs,
		cl,
		rect,
	))
	split.SetOffset(0.25)
	cl.Hide()
	rect.Hide()
	return split
}

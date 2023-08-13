package ui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/eliiasg/editor/base/state"
)

func NewTabArea(app *state.EditorApp) fyne.CanvasObject {
	tabs := container.NewDocTabs()
	tabMan := app.TabManager()
	update := func() {
		tabs.SetItems(getTabs(app))
	}
	tabs.CloseIntercept = func(ti *container.TabItem) {
		idx := slices.Index(tabs.Items, ti)
		tabMan.Tabs()[idx].CloseButtonPressed()
	}
	tabMan.UpdateUI = update
	update()
	return tabs
}

func getTabs(app *state.EditorApp) []*container.TabItem {
	tabs := app.TabManager().Tabs()
	tabItems := make([]*container.TabItem, len(tabs))
	for i, tab := range tabs {
		item := container.NewTabItem(tab.Title(), tab.Content())
		tabItems[i] = item
	}
	return tabItems
}

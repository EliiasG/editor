package state

import "fyne.io/fyne/v2"

type TabManager struct {
	tabs             []Tab
	selectedTabindex int
	UpdateUI         func()
}

func NewTabManager() *TabManager {
	return &TabManager{}
}

func (t *TabManager) SetSelected(val int) {
	t.selectedTabindex = val
}

func (t *TabManager) SelectedIndex() int {
	return t.selectedTabindex
}

func (t *TabManager) SelectedTab() Tab {
	return t.tabs[t.selectedTabindex]
}

func (t *TabManager) Tabs() []Tab {
	return t.tabs
}

func (t *TabManager) SwapTabs(idx1, idx2 int) {
	//funny swap
	t.tabs[idx1], t.tabs[idx2] = t.tabs[idx2], t.tabs[idx1]

}

type Tab interface {
	Title() string
	Content() fyne.CanvasObject
	CloseButtonPressed() bool
	ShortcutPressed(name string)
}

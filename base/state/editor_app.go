package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/x/fyne/theme"
)

type EditorApp struct {
	projectManager *ProjectManager
	tabManager     *TabManager
	app            fyne.App
	mainWindow     fyne.Window
	fileHandlers   map[string]func(string)
}

func (e *EditorApp) AddFileExtention(extentionWithoutDot string, open func(string)) {
	e.fileHandlers[extentionWithoutDot] = open
}

func (e *EditorApp) AddShortcut(name string, shortcut fyne.Shortcut) {
	e.mainWindow.Canvas().AddShortcut(shortcut, func(_ fyne.Shortcut) {
		e.tabManager.SelectedTab().ShortcutPressed(name)
	})
}

func (e *EditorApp) TabManager() *TabManager {
	return e.tabManager
}

func (e *EditorApp) AddFeature(feature Feature) {
	feature.Init(e)
}

func (e *EditorApp) ProjectManager() *ProjectManager {
	return e.projectManager
}

func (e *EditorApp) Run(root fyne.CanvasObject) {
	e.mainWindow.Resize(fyne.NewSize(800, 450))
	e.mainWindow.SetContent(root)
	e.mainWindow.ShowAndRun()
}

func (e *EditorApp) MainWindow() fyne.Window {
	return e.mainWindow
}

func NewEditorApp() EditorApp {
	eApp := EditorApp{}
	eApp.app = app.NewWithID("eliiasg.editor")
	eApp.fileHandlers = make(map[string]func(string))
	eApp.app.Settings().SetTheme(theme.AdwaitaTheme())
	eApp.projectManager = LoadProjectManagerFromSave(&eApp)
	eApp.tabManager = NewTabManager()
	eApp.mainWindow = eApp.app.NewWindow("Editor")
	eApp.mainWindow.SetMaster()
	return eApp
}

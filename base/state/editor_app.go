package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/x/fyne/theme"
)

type EditorApp struct {
	projectManager *ProjectManager
	app            fyne.App
	mainWindow     fyne.Window
}

func (e *EditorApp) AddFileExtention(extentionWithoutDot string, open func(string)) {
	//TODO implement file extentiofyne
	panic("not implemented")
}

func (e *EditorApp) OpenTab(tab Tab) {
	//TODO implement opening tabs
	panic("not implemented")
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
	eApp.app.Settings().SetTheme(theme.AdwaitaTheme())
	eApp.projectManager = LoadProjectManagerFromSave(&eApp)
	eApp.mainWindow = eApp.app.NewWindow("Editor")
	eApp.mainWindow.SetMaster()
	return eApp
}

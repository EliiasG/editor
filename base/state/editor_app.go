package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type EditorApp struct {
	projectManager *ProjectManager
	app            fyne.App
	mainWindow     fyne.Window
}

func (e *EditorApp) AddFileExtention(open func(string)) {
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

func (e *EditorApp) GetProjectManager() *ProjectManager {
	return e.projectManager
}

func (e *EditorApp) Run(root fyne.CanvasObject) {
	e.mainWindow.Resize(fyne.NewSize(800, 450))
	e.mainWindow.SetContent(root)
	e.mainWindow.ShowAndRun()
}

func (e *EditorApp) GetMainWindow() fyne.Window {
	return e.mainWindow
}

func NewEditorApp() EditorApp {
	eApp := EditorApp{}
	eApp.app = app.NewWithID("eliiasg.editor")
	eApp.projectManager = LoadProjectManagerFromSave(&eApp)
	eApp.mainWindow = eApp.app.NewWindow("Editor")
	eApp.mainWindow.SetMaster()
	return eApp
}

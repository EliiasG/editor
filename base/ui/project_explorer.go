package ui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/fileactions"
	"github.com/eliiasg/editor/base/state"
	"github.com/eliiasg/editor/base/ui/filetree"
)

func NewProjectExplorer(app *state.EditorApp) fyne.CanvasObject {
	//exp := ProjectExplorer{}
	//exp.ExtendBaseWidget(exp)
	projMan := app.ProjectManager()
	tree := filetree.NewFileTree(storage.NewFileURI(projMan.Path()))
	tree.RightClicked = func(id widget.TreeNodeID, e *fyne.PointEvent) {
		showProjectManagerContextMenu(storage.NewFileURI(id), projMan.FileActions(), e.AbsolutePosition, app.MainWindow())
	}
	return tree
}

func showProjectManagerContextMenu(uri fyne.URI, actions fileactions.FileActions, position fyne.Position, window fyne.Window) {
	menu := newDirectoryMenu(uri, actions, window)

	popUpMenu := widget.NewPopUpMenu(menu, window.Canvas())

	popUpMenu.ShowAtPosition(position)
}

func newDirectoryMenu(uri fyne.URI, actions fileactions.FileActions, window fyne.Window) *fyne.Menu {
	println(getFSPath(uri))
	info, _ := os.Stat(getFSPath(uri))
	if info.IsDir() {
		return fyne.NewMenu(
			"",
			getRenameItem(uri, actions, window),
			getDeleteItem(uri, actions, window, true),
			fyne.NewMenuItemSeparator(),
			getCutItem(uri, actions),
			getCopyItem(uri, actions),
			getPasteItem(uri, actions),
			fyne.NewMenuItemSeparator(),
			getMakeDirItem(uri, actions, window),
			getMakeFileItem(uri, actions, window),
			fyne.NewMenuItemSeparator(),
			getRevealItem(uri, actions),
		)
	} else {
		return fyne.NewMenu(
			"",
			getRenameItem(uri, actions, window),
			getDeleteItem(uri, actions, window, true),
			fyne.NewMenuItemSeparator(),
			getCutItem(uri, actions),
			getCopyItem(uri, actions),
			fyne.NewMenuItemSeparator(),
			getRevealItem(uri, actions),
		)
	}
}

func getText(title, prompt, initial string, confirmed func(string), window fyne.Window) {
	entry := widget.NewEntry()
	label := widget.NewLabel(prompt)
	label.Alignment = fyne.TextAlignCenter
	entry.SetText(initial)
	entry.OnSubmitted = func(s string) {
		confirmed(entry.Text)
	}
	menu := dialog.NewCustomConfirm(
		title,
		"Confirm",
		"Cancel",
		container.NewVBox(
			label,
			entry,
		),
		func(b bool) {
			if b {
				confirmed(entry.Text)
			}
		},
		window,
	)

	menu.Show()
	window.Canvas().Focus(entry)
	// hacky way to select the text
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})
}

func confirm(title, prompt string, confirmed func(), window fyne.Window) {
	dialog.NewConfirm(title, prompt, func(b bool) {
		if b {
			confirmed()
		}
	}, window).Show()
}

func getFSPath(uri fyne.URI) string {
	return uri.Path()[7:]
}

func getRenameItem(uri fyne.URI, actions fileactions.FileActions, window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Rename",
		func() {
			getText(
				"Rename",
				"Rename \""+uri.Name()+"\":",
				uri.Name(),
				func(s string) {
					actions.Rename(getFSPath(uri), s)
				},
				window,
			)
		},
	)
}

func getDeleteItem(uri fyne.URI, actions fileactions.FileActions, window fyne.Window, dir bool) *fyne.MenuItem {
	typeName := "File"
	if dir {
		typeName = "Directory"
	}
	return fyne.NewMenuItem(
		"Delete",
		func() {
			confirm(
				"Delete "+typeName,
				"Are you sure you want to delete \""+uri.Name()+"\"?",
				func() {
					actions.Delete(getFSPath(uri))
				},
				window,
			)
		},
	)
}

func getCutItem(uri fyne.URI, actions fileactions.FileActions) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Cut",
		func() {
			actions.Cut(getFSPath(uri))
		},
	)
}

func getCopyItem(uri fyne.URI, actions fileactions.FileActions) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Copy",
		func() {
			actions.Copy(getFSPath(uri))
		},
	)
}

func getPasteItem(uri fyne.URI, actions fileactions.FileActions) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Paste",
		func() {
			actions.Paste(getFSPath(uri))
		},
	)
}

func getMakeDirItem(uri fyne.URI, actions fileactions.FileActions, window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Add Directory",
		func() {
			getText(
				"Add Directory",
				"Add Directory:",
				"NewDir",
				func(s string) {
					actions.CreateFolder(getFSPath(uri) + "/" + s)
				},
				window,
			)
		},
	)
}

func getMakeFileItem(uri fyne.URI, actions fileactions.FileActions, window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Add File",
		func() {
			getText(
				"Add File",
				"Add File:",
				"NewFile",
				func(s string) {
					actions.CreateFile(getFSPath(uri) + "/" + s)
				},
				window,
			)
		},
	)
}

func getRevealItem(uri fyne.URI, actions fileactions.FileActions) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Show in explorer",
		func() {
			actions.Reveal(getFSPath(uri))
		},
	)
}

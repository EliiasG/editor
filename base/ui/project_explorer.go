package ui

import (
	"os"
	"time"

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
	tree.Refresh()
	tree.RightClicked = func(id widget.TreeNodeID, e *fyne.PointEvent) {
		showProjectManagerContextMenu(
			storage.NewFileURI(id),
			projMan.FileActions(),
			e.AbsolutePosition,
			func() {
				reload(&tree.Tree)
			},
			app.MainWindow(),
		)
	}
	handleSelect(&tree.Tree, projMan.FileActions())
	return tree
}

func handleSelect(tree *widget.Tree, actions fileactions.FileActions) {
	prev := ""
	var prevTime int64
	tree.OnSelected = func(uid widget.TreeNodeID) {
		tree.UnselectAll()
		if tree.IsBranch(uid) {
			tree.ToggleBranch(uid)
			return
		}
		now := time.Now().UnixMilli()
		if prev == uid && prevTime+DoubleClickTimeMS >= now {
			actions.Open(getFSPath(storage.NewFileURI(uid)))
		}
		prev = uid
		prevTime = now
	}
}

func reload(tree *widget.Tree) {
	reloadUID(tree, tree.Root)
}

func reloadUID(tree *widget.Tree, uid widget.TreeNodeID) {
	if !tree.IsBranchOpen(uid) {
		return
	}
	// whacky hacky to get node to reload
	tree.CloseBranch(uid)
	tree.OpenBranch(uid)
	for _, uid := range tree.ChildUIDs(uid) {
		reloadUID(tree, uid)
	}
}

func showProjectManagerContextMenu(uri fyne.URI, actions fileactions.FileActions, position fyne.Position, reload func(), window fyne.Window) {
	menu := newFSItemMenu(uri, actions, reload, window)

	popUpMenu := widget.NewPopUpMenu(menu, window.Canvas())

	popUpMenu.ShowAtPosition(position)
}

func newFSItemMenu(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window) *fyne.Menu {
	info, _ := os.Stat(getFSPath(uri))
	if info.IsDir() {
		return fyne.NewMenu(
			"",
			getRenameItem(uri, actions, reload, window),
			getDeleteItem(uri, actions, reload, window, true),
			fyne.NewMenuItemSeparator(),
			getCutItem(uri, actions, reload),
			getCopyItem(uri, actions, reload),
			getPasteItem(uri, actions, reload),
			fyne.NewMenuItemSeparator(),
			getMakeDirItem(uri, actions, reload, window),
			getMakeFileItem(uri, actions, reload, window),
			fyne.NewMenuItemSeparator(),
			getRevealItem(uri, actions, reload),
		)
	} else {
		return fyne.NewMenu(
			"",
			getRenameItem(uri, actions, reload, window),
			getDeleteItem(uri, actions, reload, window, false),
			fyne.NewMenuItemSeparator(),
			getCutItem(uri, actions, reload),
			getCopyItem(uri, actions, reload),
			fyne.NewMenuItemSeparator(),
			getRevealItem(uri, actions, reload),
		)
	}
}

func getText(title, prompt, initial string, confirmed func(string), window fyne.Window) {
	entry := widget.NewEntry()
	label := widget.NewLabel(prompt)
	label.Alignment = fyne.TextAlignCenter
	entry.SetText(initial)
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

	entry.OnSubmitted = func(s string) {
		confirmed(entry.Text)
		menu.Hide()
	}

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

func getRenameItem(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Rename",
		func() {
			getText(
				"Rename",
				"Rename \""+uri.Name()+"\":",
				uri.Name(),
				func(s string) {
					actions.Rename(getFSPath(uri), s)
					reload()
				},
				window,
			)
		},
	)
}

func getDeleteItem(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window, dir bool) *fyne.MenuItem {
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
					reload()
				},
				window,
			)
		},
	)
}

func getCutItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Cut",
		func() {
			actions.Cut(getFSPath(uri))
			reload()
		},
	)
}

func getCopyItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Copy",
		func() {
			actions.Copy(getFSPath(uri))
			reload()
		},
	)
}

func getPasteItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	item := fyne.NewMenuItem(
		"Paste",
		func() {
			actions.Paste(getFSPath(uri))
			reload()
		},
	)
	item.Disabled = !actions.HasClipboardContent()
	return item
}

func getMakeDirItem(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Add Directory",
		func() {
			getText(
				"Add Directory",
				"Add Directory:",
				"NewDir",
				func(s string) {
					actions.CreateFolder(getFSPath(uri) + "/" + s)
					reload()
				},
				window,
			)
		},
	)
}

func getMakeFileItem(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Add File",
		func() {
			getText(
				"Add File",
				"Add File:",
				"NewFile",
				func(s string) {
					actions.CreateFile(getFSPath(uri) + "/" + s)
					reload()
				},
				window,
			)
		},
	)
}

func getRevealItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Show in explorer",
		func() {
			actions.Reveal(getFSPath(uri))
			reload()
		},
	)
}

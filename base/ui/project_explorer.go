package ui

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/eliiasg/editor/base/fileactions"
	"github.com/eliiasg/editor/base/fileutil"
	"github.com/eliiasg/editor/base/state"
	"github.com/eliiasg/editor/base/ui/filetree"
)

// this file contains some of the worst ui code to ever exist, please read at your own risk

func NewProjectExplorer(app *state.EditorApp, click func(string)) (fyne.CanvasObject, func()) {
	projMan := app.ProjectManager()
	tree := filetree.NewFileTree(storage.NewFileURI(projMan.Path()))
	filter := fileutil.NewSearchFilter(projMan.Path())
	tree.Filter = filter
	tree.RightClicked = func(id widget.TreeNodeID, e *fyne.PointEvent) {
		showProjectManagerContextMenu(
			storage.NewFileURI(id),
			projMan.FileActions(),
			e.AbsolutePosition,
			func() {
				reload(&tree.Tree)
				filter.UpdateCache()
			},
			app.MainWindow(),
		)
	}
	rl := func() {
		filter.UpdateSearch("", app.ProjectManager().FilterExtention(), true)
		reload(&tree.Tree)
	}
	handleSelect(&tree.Tree, projMan.FileActions(), func(s string) {
		click(s)
		reload(&tree.Tree)
	})
	return container.NewBorder(
		newEntrySection(filter, app, func() {
			go reload(&tree.Tree)
		}, &tree.Tree),
		nil,
		nil,
		nil,
		tree,
	), rl
}

func newEntrySection(filter *fileutil.SearchFilter, app *state.EditorApp, changed func(), tree *widget.Tree) fyne.CanvasObject {
	entry := widget.NewEntry()
	shortcut := &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierAlt}
	canvas := app.MainWindow().Canvas()
	canvas.AddShortcut(shortcut, func(shortcut fyne.Shortcut) {
		canvas.Focus(entry)
		entry.TypedShortcut(&fyne.ShortcutSelectAll{})
	})
	clearButton := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		entry.SetText("")
	})
	var state []string
	var old string
	entry.OnChanged = func(s string) {
		if old == "" {
			state = getTreeState(tree)
		}
		if s == "" {
			setTreeState(tree, state)
		} else {
			tree.OpenAllBranches()
		}
		println(app.ProjectManager().FilterExtention())
		go filter.UpdateSearch(s, app.ProjectManager().FilterExtention(), true)
		changed()
		old = s
	}
	return container.NewBorder(
		nil,
		nil,
		nil,
		clearButton,
		entry,
	)
}

func setTreeState(tree *widget.Tree, state []string) {
	tree.CloseAllBranches()
	for _, node := range state {
		tree.OpenBranch(node)
	}
}

func getTreeState(tree *widget.Tree) []string {
	res := make([]string, 0)
	getTreeUIDStates(tree, tree.Root, &res)
	return res
}

func getTreeUIDStates(tree *widget.Tree, uid widget.TreeNodeID, state *[]string) {
	if !tree.IsBranchOpen(uid) {
		return
	}
	*state = append(*state, uid)
	for _, uid := range tree.ChildUIDs(uid) {
		getTreeUIDStates(tree, uid, state)
	}
}

func handleSelect(tree *widget.Tree, actions fileactions.FileActions, click func(string)) {
	prev := ""
	var prevTime int64
	tree.OnSelected = func(uid widget.TreeNodeID) {
		tree.UnselectAll()
		if tree.IsBranch(uid) {
			tree.ToggleBranch(uid)
			return
		}
		path := fileutil.GetFSPath(storage.NewFileURI(uid))
		click(path)
		now := time.Now().UnixMilli()
		if prev == uid && prevTime+DoubleClickTimeMS >= now {
			actions.Open(path)
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
	// whacky hacky way to get node to reload
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
	info, _ := os.Stat(fileutil.GetFSPath(uri))
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

func getRenameItem(uri fyne.URI, actions fileactions.FileActions, reload func(), window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Rename",
		func() {
			getText(
				"Rename",
				"Rename \""+uri.Name()+"\":",
				uri.Name(),
				func(s string) {
					actions.Rename(fileutil.GetFSPath(uri), s)
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
					actions.Delete(fileutil.GetFSPath(uri))
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
			actions.Cut(fileutil.GetFSPath(uri))
			reload()
		},
	)
}

func getCopyItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	return fyne.NewMenuItem(
		"Copy",
		func() {
			actions.Copy(fileutil.GetFSPath(uri))
			reload()
		},
	)
}

func getPasteItem(uri fyne.URI, actions fileactions.FileActions, reload func()) *fyne.MenuItem {
	item := fyne.NewMenuItem(
		"Paste",
		func() {
			actions.Paste(fileutil.GetFSPath(uri))
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
					actions.CreateFolder(fileutil.GetFSPath(uri) + "/" + s)
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
					actions.CreateFile(fileutil.GetFSPath(uri) + "/" + s)
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
			actions.Reveal(fileutil.GetFSPath(uri))
			reload()
		},
	)
}

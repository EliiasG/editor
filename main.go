package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	/*
		hello := widget.NewLabel("Hello Fyne!")
		w.SetContent(container.NewVBox(
			hello,
			widget.NewButton("Hi!", func() {
				hello.SetText("Welcome :)")
			}),
		))
	*/
	tree := xwidget.NewFileTree(storage.NewFileURI("C:\\Projects"))    // Start from home directory
	tree.Filter = storage.NewExtensionFileFilter([]string{".txt", ""}) // Filter files
	tree.Sorter = func(u1, u2 fyne.URI) bool {
		return u1.String() < u2.String() // Sort alphabetically
	}
	tree.OnSelected = func(uid widget.TreeNodeID) {
		if tree.IsBranchOpen(uid) {
			tree.CloseBranch(uid)
		} else {
			tree.OpenBranch(uid)
		}
		tree.Unselect(uid)
		println(uid)
	}
	w.SetContent(tree)
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}

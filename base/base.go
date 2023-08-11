package base

import (
	"github.com/eliiasg/editor/base/state"
	"github.com/eliiasg/editor/base/ui"
)

func RunApp(features []state.Feature) {
	app := state.NewEditorApp()
	for _, v := range features {
		app.AddFeature(v)
	}
	app.Run(ui.GetRoot(&app))
}

package base

type EditorApp struct {
}

func (e *EditorApp) AddExtention(open func(string)) {
	panic("not implemented")
}

func (e *EditorApp) OpenTab() {
	panic("not implemented")
}

func (e *EditorApp) AddFeature(feature Feature) {
	feature.Init(e)
}

package fileactions

type FileActions interface {
	Open(path string)
	Delete(path string)
	CreateFile(path string)
	CreateFolder(path string)
	Cut(path string)
	Copy(path string)
	Paste(path string)
	Rename(path, name string)
	Reveal(path string)
	HasClipboardContent() bool
}

type simpleImpl struct{}

func (f simpleImpl) Open(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Delete(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) CreateFile(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) CreateFolder(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Cut(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Copy(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Paste(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Rename(path, name string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) Reveal(path string) {
	panic("not implemented") // TODO: Implement
}

func (f simpleImpl) HasClipboardContent() bool {
	return false
	//TODO
}

func NewSimple() FileActions {
	return simpleImpl{}
}

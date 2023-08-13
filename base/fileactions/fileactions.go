package fileactions

import (
	"os"
	pth "path"

	cp "github.com/otiai10/copy"
	"github.com/skratchdot/open-golang/open"
)

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

type simpleImpl struct {
	copied  string
	cutting bool
	open    func(string)
}

func (f *simpleImpl) Open(path string) {
	f.open(path)
}

func (f *simpleImpl) Delete(path string) {
	os.RemoveAll(path)
}

func (f *simpleImpl) CreateFile(path string) {
	file, _ := os.Create(path)
	file.Close()
}

func (f *simpleImpl) CreateFolder(path string) {
	os.MkdirAll(path, 0777)
}

func (f *simpleImpl) Cut(path string) {
	f.copied = path
	f.cutting = true
}

func (f *simpleImpl) Copy(path string) {
	f.copied = path
	f.cutting = false
}

func (f *simpleImpl) Paste(path string) {
	newPath := pth.Join(path, pth.Base(f.copied))
	for exists(newPath) {
		newPath += "_"
	}
	err := cp.Copy(f.copied, newPath)
	if err == nil && f.cutting {
		os.RemoveAll(f.copied)
	}
}

func (f *simpleImpl) Rename(path, name string) {
	os.Rename(path, pth.Join(pth.Dir(path), name))
}

func (f *simpleImpl) Reveal(path string) {
	open.Start(pth.Dir(path))
}

func (f *simpleImpl) HasClipboardContent() bool {
	return f.copied != ""
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func NewSimple(open func(string)) FileActions {
	return &simpleImpl{open: open}
}

package fileutil

import (
	"os"
	pth "path"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
)

func GetFSPath(uri fyne.URI) string {
	return uri.Path()[7:]
}

type SearchFilter struct {
	//[path] child paths
	cache    map[string][]string
	matching map[string]bool
	root     string
	mu       sync.Mutex
}

func NewSearchFilter(root string) *SearchFilter {
	filter := &SearchFilter{
		cache:    make(map[string][]string),
		matching: make(map[string]bool),
		root:     root,
	}
	filter.UpdateCache()
	filter.UpdateSearch("", "", true)
	return filter
}

func (f *SearchFilter) Matches(uri fyne.URI) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.matching[uri.Path()]
}

func (f *SearchFilter) UpdateSearch(search, extentionWithoutDot string, ignoreCase bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.matching {
		delete(f.matching, k)
	}
	f.addMatching(f.root, search, extentionWithoutDot, ignoreCase)
	// always show root
	f.matching[f.root] = true
}

func (f *SearchFilter) addMatching(path, search, extension string, ignoreCase bool) bool {
	// matches extention or search
	matching := true
	sPath := path
	if ignoreCase {
		sPath = strings.ToLower(path)
		search = strings.ToLower(search)
	}
	if extension != "" {
		ext := pth.Ext(path)
		matching = ext != "" && ext[1:] == extension
	}
	if search != "" {
		matching = matching && strings.Contains(sPath, search)
	}
	children := f.cache[path]
	// add matching for children, and set self if any child matches
	if children != nil {
		for _, child := range children {
			if f.addMatching(child, search, extension, ignoreCase) {
				matching = true
			}
		}
	}
	f.matching[path] = matching
	return matching
}

func (f *SearchFilter) UpdateCache() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.cache {
		delete(f.cache, k)
	}
	f.addToCache(f.root)
}

func (f *SearchFilter) addToCache(path string) {
	inf, err := os.Stat(path)
	if err != nil {
		return
	}
	if !inf.IsDir() {
		// no files as keys
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		// probably missing permission or something
		return
	}
	res := make([]string, len(entries))
	for i, entry := range entries {
		entryPath := pth.Join(path, entry.Name())
		res[i] = entryPath
		if entry.IsDir() {
			defer f.addToCache(entryPath)
		}
	}
	f.cache[path] = res
}

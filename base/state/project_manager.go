package state

import (
	pth "path"
	"slices"
	"strings"

	"github.com/eliiasg/editor/base/fileactions"
	op "github.com/skratchdot/open-golang/open"
)

type ProjectManager struct {
	isOpen         bool
	openPath       string
	recentProjects []string
	fileActions    fileactions.FileActions
}

func (p *ProjectManager) IsOpen() bool {
	return p.isOpen
}

func (p *ProjectManager) Path() string {
	return p.openPath
}

func (p *ProjectManager) AddProject(path string) {
	if !slices.Contains(p.recentProjects, path) {
		p.recentProjects = append([]string{path}, p.recentProjects...)
	}
}

func (p *ProjectManager) RemoveProject(idx int) {
	p.recentProjects = append(p.recentProjects[0:idx], p.recentProjects[idx+1:]...)
}

func (p *ProjectManager) Save(app *EditorApp) {
	app.app.Preferences().SetString("recent", strings.Join(p.recentProjects, "\n"))
}

func (p *ProjectManager) OpenProject(idx int) {
	p.openPath = p.recentProjects[idx]
	p.isOpen = true
	// to move to top
	p.RemoveProject(idx)
	p.AddProject(p.openPath)
}

func (p *ProjectManager) RecentProjects() []string {
	return p.recentProjects
}

func (p *ProjectManager) FileActions() fileactions.FileActions {
	return p.fileActions
}

func LoadProjectManagerFromSave(app *EditorApp) *ProjectManager {
	open := func(path string) {
		ext := pth.Ext(path)[1:]
		if handler := app.fileHandlers[ext]; handler == nil {
			op.Start(path)
		} else {
			handler(path)
		}
	}

	projects := app.app.Preferences().StringWithFallback("recent", "")
	var projectSlice []string
	if projects != "" {
		projectSlice = strings.Split(projects, "\n")
	}
	return &ProjectManager{
		recentProjects: projectSlice,
		fileActions:    fileactions.NewSimple(open),
	}
}

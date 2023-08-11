package state

import "strings"

type ProjectManager struct {
	isOpen         bool
	openPath       string
	recentProjects []string
}

func (p *ProjectManager) IsOpen() bool {
	return p.isOpen
}

func (p *ProjectManager) GetPath() string {
	return p.openPath
}

func (p *ProjectManager) AddProject(path string) {
	p.recentProjects = append([]string{path}, p.recentProjects...)
}

func (p *ProjectManager) RemoveProject(idx int) {
	p.recentProjects = append(p.recentProjects[0:idx], p.recentProjects[idx+1:]...)
}

func (p *ProjectManager) Save(app *EditorApp) {
	app.app.Preferences().SetString("recent", strings.Join(p.recentProjects, "\n"))
}

func (p *ProjectManager) Open(idx int) {
	p.openPath = p.recentProjects[idx]
	p.isOpen = true
	// to move to top
	p.RemoveProject(idx)
	p.AddProject(p.openPath)
}

func (p *ProjectManager) GetRecentProjects() []string {
	return p.recentProjects
}

func LoadProjectManagerFromSave(app *EditorApp) *ProjectManager {
	projects := app.app.Preferences().StringWithFallback("recent", "")
	var projectSlice []string
	if projects != "" {
		projectSlice = strings.Split(projects, "\n")
	}
	return &ProjectManager{
		recentProjects: projectSlice,
	}
}

package base

type ProjectManager struct {
	is_open         bool
	open_path       string
	recent_projects []string
}

func (p *ProjectManager) IsOpen() bool {
	return p.is_open
}

func (p *ProjectManager) GetPath() string {
	return p.open_path
}

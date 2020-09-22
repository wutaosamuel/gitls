package gitls

// GitBranch contain git branch infomation
type GitBranch struct {
	Name string		// Name   -> name for local branch
	Status string // Status -> branch status against remote
}

// NewGitBranch create a new GitBranch
func NewGitBranch() *GitBranch {
	return &GitBranch{
		Name: "",
		Status: ""}
}
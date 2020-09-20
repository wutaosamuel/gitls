package gitls

// GitBranch contain git branch infomation
type GitBranch struct {
	Name string		// Name -> name for local branch
	Remote string // Remote -> upstream remote name
}
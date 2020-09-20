package gitls

// GitRemote contain remote url information for pushing & fetching
type GitRemote struct {
	Name  string   // Name  -> git remote name
	Fetch string   // Fetch -> url fetch from remote
	Push  []string // Push  -> urls for pushing to remote
}

// NewGitRemote create a new GitRemote object
func NewGitRemote() *GitRemote {
	return &GitRemote{
		Name:  "",
		Fetch: "",
		Push:  make([]string, 0)}
}

// SetGitRemote set information
func (g *GitRemote) SetGitRemote(name, fetch string, push []string) {
	g.Name = name
	g.Fetch = fetch
	g.Push = push
}

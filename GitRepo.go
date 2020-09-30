package gitls

import (
	"bytes"
	"sync"

	gogit "github.com/go-git/go-git/v5"
	// gogit "github.com/go-git/go-git"
)

// GitRepo contain infomation of remotes && branches
type GitRepo struct {
	Repo     *gogit.Repository       // Repo 		 -> git repository
	Dir      string                  // Dir 		 -> git dir
	Remotes  map[string]GitRemote    // Remotes  -> git remotes infomation
	Branches map[string][]*GitBranch // Branches -> git branches infomation

	rw *sync.RWMutex
}

// NewGitRepo create a new GitRepo
/*
 *	- check git is exist in system
 */
func NewGitRepo() *GitRepo {
	if !CheckGitExist() {
		panic("This version requires git. Download git program first, please!")
	}

	return &GitRepo{
		Repo:     &gogit.Repository{},
		Dir:      "",
		Remotes:  make(map[string]GitRemote, 0),
		Branches: make(map[string][]*GitBranch, 0),
		rw:       &sync.RWMutex{}}
}

// CheckGitDir check whether dir is git dir or not
/*
 *	- if yes, set repo
 *  - if no, return false
 */
func (g *GitRepo) CheckGitDir(dir string) (bool, error) {
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		if err.Error() == "repository does not exist" {
			return false, nil
		}
		return false, err
	}
	g.Repo = repo
	g.Dir = dir

	return true, nil
}

// GetGitInfo get information of git directory
func (g *GitRepo) GetGitInfo(dir string) error {
	g.rw.Lock()
	defer g.rw.Unlock()

	// check git dir && set git repo
	_, err := g.CheckGitDir(dir)
	if err != nil {
		return err
	}

	// set remotes
	remotes, err := g.Repo.Remotes()
	if err != nil {
		return err
	}
	for _, r := range remotes {
		remoteConfig := r.Config()
		remote := g.Remotes[remoteConfig.Name]
		for k, url := range remoteConfig.URLs {
			if k == 0 {
				remote.Fetch = url
			}
			remote.Push = append(remote.Push, url)
		}
		g.Remotes[remoteConfig.Name] = remote
	}

	// set branch
	for k := range g.Remotes {
		branches := g.Branches[k]
		info, _ := RunGitShow(dir, k)
		for _, m := range info["push"] {
			branch := NewGitBranch()
			for b, s := range GetPushStatus(m) {
				branch.Name = b
				branch.Status = s
			}
			branches = append(branches, branch)
		}
		g.Branches[k] = branches
	}

	return nil
}

// DefaultString console output by default
/*
 *	GitDirectory
 *	RemoteName
 *		branch1		status
 *		branch2		status
 */
func (g *GitRepo) DefaultString() string {
	buffer := bytes.NewBuffer(make([]byte, 0, 3))

	// Git name
	buffer.WriteString(g.Dir)
	buffer.WriteByte(':')
	buffer.WriteByte('\n')
	// RemoteName && branchs && status
	buffer.Write(g.DefaultElementString())

	return buffer.String()
}

// URLWithDirString display dir && urlstring
/*
 *	GitDir
 *	RemoteName
 *		Fetch: url1
 *		Pull:  url1
 *					 url2
 */
func (g *GitRepo) URLWithDirString() string {
	buffer := &bytes.Buffer{}

	buffer.WriteString(g.Dir)
	buffer.WriteByte('\n')
	for name, remote := range g.Remotes {
		buffer.WriteString(name)
		buffer.WriteByte('\n')
		buffer.Write(g.URLString(remote.Fetch, remote.Push))
	}

	return buffer.String()
}

// DefaultElementString element string
/*
 *	RemoteName
 *		branch1		status
 *		branch2		status
 */
func (g *GitRepo) DefaultElementString() []byte {
	buffer := &bytes.Buffer{}

	for remote, branches := range g.Branches {
		buffer.WriteString(remote)
		buffer.WriteByte('\n')
		for _, branch := range branches {
			buffer.Write(g.BranchStatusString(branch.Name, branch.Status))
		}
		buffer.WriteByte('\n')
	}

	return buffer.Bytes()
}

// AllString display all infomation
/*
 *	GitDir
 *	RemoteName
 *		Fetch: url1
 *		push:  url1
 *		       url2
 *	Local branches against remote branch with status
 *		branch1		status
 *		branch2		status
 */
func (g *GitRepo) AllString() string {
	buffer := bytes.NewBuffer(make([]byte, 0, 1))

	// GitDir
	buffer.WriteString(g.Dir)
	buffer.WriteByte('\n')
	// Remote
	for name, remote := range g.Remotes {
		// RemoteName
		buffer.WriteString(name)
		buffer.WriteByte('\n')
		buffer.Write(g.URLString(remote.Fetch, remote.Push))
		// buffer.WriteByte('\n')
		// buffer.WriteString("Local branches against remote branch with status")
		// buffer.WriteByte('\n')
		branches := g.Branches[name]
		for _, b := range branches {
			buffer.Write(g.BranchStatusString(b.Name, b.Status))
		}
		buffer.WriteByte('\n')
	}

	return buffer.String()
}

// BranchStatusString display branch && relate status
/*
 *	branch1		status
 */
func (g *GitRepo) BranchStatusString(branch, status string) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 5))

	buffer.WriteString(branch)
	buffer.WriteByte('\t')
	buffer.WriteString("->")
	buffer.WriteByte('\t')
	buffer.WriteString(status)
	buffer.WriteByte('\n')

	return buffer.Bytes()
}

// URLString display all urls for fetch && pull
/*
 *	URLString
 *
 *	Fetch: url
 *	Push: url
 *				url2
 */
func (g *GitRepo) URLString(fetch string, pushes []string) []byte {
	buffer := &bytes.Buffer{}

	if fetch != "" {
		buffer.WriteString("Fetch:")
		buffer.WriteByte('\t')
		buffer.WriteString(fetch)
		buffer.WriteByte('\n')
	}
	if len(pushes) != 0 {
		for k, v := range pushes {
			if k == 0 {
				buffer.WriteString("Pull:")
				buffer.WriteByte('\t')
				buffer.WriteString(v)
				buffer.WriteByte('\n')
			}
			if k != 0 {
				buffer.WriteByte('\t')
				buffer.WriteString(v)
				buffer.WriteByte('\n')
			}
		}
	}

	return buffer.Bytes()
}

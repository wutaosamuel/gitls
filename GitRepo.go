package gitls

import (
	"sync"

	// gogit "github.com/go-git/go-git/v5"
	gogit "github.com/go-git/go-git"
)

// GitRepo contain infomation of remotes && branches
type GitRepo struct {
	Repo     *gogit.Repository    // Repo 		-> git repository
	Remotes  map[string]GitRemote // Remotes 	-> git remotes infomation
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
		return false, err
	}
	g.Repo = repo

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

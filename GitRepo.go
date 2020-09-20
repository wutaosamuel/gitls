package gitls

import (
	// gogit "github.com/go-git/go-git/v5"
	gogit "github.com/go-git/go-git"
)

// GitRepo contain infomation of remotes && branches
type GitRepo struct {
	Repo  *gogit.Repository  // Repo -> git repository
	Remotes map[string]GitRemote // Remotes -> git remotes infomation
	// Branches
}
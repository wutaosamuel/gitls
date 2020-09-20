// +build !windows

package gitls

import (
	"os/exec"
)

// RunGitShow run git command
//	- need to check git first, TODO:
func RunGitShow(dir, remote string) (map[string][]string, error) {
	command := []string{"git", "remote", "show", remote}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return cutMessage(string(output)), nil
}

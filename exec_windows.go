// +build windows

package gitls

import (
	"os"
	"os/exec"
)

// RunGitShow run git command
//	- need to check git first, TODO:
func RunGitShow(dir, remote string) (map[string][]string, error) {
	command := []string{"git", "remote", "show", remote}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return cutMessage(string(output)), nil
}

// CheckGitExist check git can be execute
func CheckGitExist() bool {
	cmd := exec.Command("git", "--version")
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
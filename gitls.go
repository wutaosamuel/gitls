package gitls

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/pflag"
)

var (
	helpFlag   = pflag.BoolP("help", "h", false, "Display Usage")
	inputFlag  = pflag.StringP("input", "i", "", "The git folder (single)")
	inputsFlag = pflag.StringSliceP("inputs", "I", []string{}, "The gits' folder(multiple)")
	allFlag    = pflag.BoolP("all", "a", false, "The git info. (URL, remote, status)")
	urlFlag    = pflag.BoolP("url", "u", false, "Display Fetch url and Push urls")
)

// Main func of gitls
func Main() {
	var dirs []string
	var wg sync.WaitGroup

	pflag.Parse()
	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	// handle input
	// handle single input
	if *inputFlag != "" {
		// check dir is exist && get absolute path
		dir, err := filepath.Abs(*inputFlag)
		if err != nil {
			printErrorUsage(err)
		}
		isDir, err := IsDir(dir)
		if err != nil {
			printErrorUsage(err)
		}
		if !isDir {
			fmt.Println(dir + " is not a directory")
			fmt.Println()
			pflag.Usage()
			os.Exit(0)
		}
		fmt.Println(dir)
		dirs = getDirs(dir)
	}
	// handle inputs
	if len(*inputsFlag) != 0 {
		for _, v := range *inputsFlag {
			dir, err := filepath.Abs(v)
			if err != nil {
				fmt.Print(dir + ": ")
				fmt.Println(err)
			}
			isDir, err := IsDir(dir)
			if err != nil {
				fmt.Print(dir + ": ")
				fmt.Println(err)
			}
			if isDir {
				dirs = append(dirs, dir)
			}
		}
	}
	// default input
	if *inputFlag == "" && len(*inputsFlag) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			printErrorUsage(err)
		}
		dirs = getDirs(dir)
	}

	// get git dirs
	gits, dir := separateDirs(dirs)
	// get git info
	fetchInfo := func(g *GitRepo) {
		defer wg.Done()
		g.GetGitInfo(g.Dir)
	}
	for _, g := range gits {
		wg.Add(1)
		go fetchInfo(g)
	}
	wg.Wait()

	// display git infomation by flag
	// display -a | --all
	if *allFlag {
		for _, g := range gits {
			if len(g.Remotes) == 0 {
				fmt.Println(g.Dir, ":")
				fmt.Println("-> ", "no upstream remote")
				fmt.Println()
				continue
			}
			fmt.Print(g.AllString())
		}
	}
	// display -u | --url
	if *urlFlag {
		for _, g := range gits {
			if len(g.Remotes) == 0 {
				fmt.Println(g.Dir, ":")
				fmt.Println("-> ", "no upstream remote")
				fmt.Println()
				continue
			}
			fmt.Print(g.URLWithDirString())
		}
	}
	// display default output
	if !*allFlag && !*urlFlag {
		for _, g := range gits {
			if len(g.Remotes) == 0 {
				fmt.Println(g.Dir, ":")
				fmt.Println("-> ", "no upstream remote")
				fmt.Println()
				continue
			}
			fmt.Print(g.DefaultString())
		}
	}

	// display dir infomation if No. of dir > 0
	if len(dir) > 0 {
		fmt.Println()
		fmt.Println("The following folders is not initialized by git: ")
		for _, d := range dir {
			fmt.Print("-> ")
			fmt.Println(d)
		}
	}

	return
}

// getDirs get all dirs under directory, exclude .git directory
/*
 *	- add this dir
 *  - add sub-dirs
 */
func getDirs(root string) []string {
	var dirs []string
	dirs = append(dirs, root)
	files, _ := ioutil.ReadDir(root)
	for _, f := range files {
		name := f.Name()
		dir := path.Join(root, f.Name())
		isDir, _ := IsDir(dir)
		if name == ".git" {
			isDir = false
		}
		if isDir {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}

// printErrorUsage print Error info. && usage info.
func printErrorUsage(err error) {
	fmt.Println(err)
	fmt.Println()
	pflag.Usage()
	os.Exit(0)

	return
}

// separateDirs separate dirs to git dirs && non-git dirs
func separateDirs(dirs []string) ([]*GitRepo, []string) {
	var gits []*GitRepo
	var dir []string

	for _, d := range dirs {
		git := NewGitRepo()
		isGit, err := git.CheckGitDir(d)
		if err != nil {
			printErrorUsage(err)
		}
		if isGit {
			gits = append(gits, git)
		}
		if !isGit {
			dir = append(dir, d)
		}
	}

	return gits, dir
}

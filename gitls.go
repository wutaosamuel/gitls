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
		dirs = handleInput(*inputFlag)
	}
	// handle inputs
	if len(*inputsFlag) != 0 {
		dirs = handleInputs(*inputsFlag)
	}
	// default input
	if *inputFlag == "" && len(*inputsFlag) == 0 {
		dirs = handleNoInput()
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
		handleAllDisplay(gits)
	}
	// display -u | --url
	if *urlFlag {
		handleURLDisplay(gits)
	}
	// display default output
	if !*allFlag && !*urlFlag {
		handleDefaultDisplay(gits)
	}

	// display non-git dir infomation if No. of dir > 0
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

// handleInput handle single input
/*
 *	Output:
 *	- get input folder
 *  - get sub-folder under input folder
 */
func handleInput(input string) []string {
	// check dir is exist && get absolute path
	dir, err := filepath.Abs(input)
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
	dirs := getDirs(dir)

	return dirs
}

// handleInputs handle inputs
/*
 *	- get absolute path of directories
 *	- check it is directory
 */
func handleInputs(inputs []string) []string {
	var dirs []string
	for _, v := range inputs {
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

	return dirs
}

// handleNoInput process no input or inputs, which is default input
func handleNoInput() []string {
	dir, err := os.Getwd()
	if err != nil {
		printErrorUsage(err)
	}
	dirs := getDirs(dir)

	return dirs
}

// handleAllDisplay display remote name && pull&push url && branch -> status
func handleAllDisplay(gits []*GitRepo) {
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

// handleURLDisplay display remote name && pull&push url
func handleURLDisplay(gits []*GitRepo) {
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

// handleDefaultDisplay display remote name && branch -> status
func handleDefaultDisplay(gits []*GitRepo) {
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

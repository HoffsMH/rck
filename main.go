package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "rck",
	Short: "repository check",
	Long: `check all repos recursively from a given directory
					for their up-to-dateness with origin`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := filepath.Abs(args[0])
		var errors []error

		repos, err := RepoList(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, r := range repos {
			rs, err := CheckRepo(r)
			if err != nil {
				errors = append(errors, err)
			} else if !rs {
				fmt.Println(r)
			}
		}
		fmt.Println()
		fmt.Println("=========ERRORS========")
		for _, e := range errors {
			fmt.Println(e)
		}
	},
}

func RepoList(dir string) ([]string, error) {
	var list []string
	permissionDenied := regexp.MustCompile(`permission denied`)

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if !permissionDenied.MatchString(fmt.Sprintf("%+v", err)) {
					return err
				}
				return nil
			}

			if info.IsDir() {
				isRepo, _ := IsRepo(path)
				if isRepo {
					list = append(list, path)
				}
			}
			return nil
		})
	return list, err
}

func CheckRepo(repoDir string) (bool, error) {
	os.Chdir(repoDir)

	// check that there is at least one remote
	o, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return false, fmt.Errorf(repoDir, ": git remote failed")
	}
	remotes := len(strings.Split(string(o), "\n"))
	if remotes < 2 {
		return false, fmt.Errorf(repoDir, ": no remotes")
	}

	// fetch latest information from remotes
	_, err = exec.Command("git", "fetch").Output()
	if err != nil {
		return false, fmt.Errorf(repoDir, ": git fetch failed")
	}

	// check that the repo is up to date
	upToDate, err := IsRepoUpToDate(repoDir)
	if err != nil {
		return false, err
	}
	return upToDate, err
}

func IsRepoUpToDate(dir string) (bool, error) {
	localOnly, err := BranchHasLocalOnly(dir)
	if err != nil {
		return false, err
	}
	treeClean, err := IsTreeClean(dir)
	if err != nil {
		return false, err
	}

	return !localOnly && treeClean, nil
}

func BranchHasLocalOnly(dir string) (bool, error) {
	// get current branch name
	o, err := exec.Command("git", "rev-parse","--abbrev-ref", "HEAD").Output()
	bn := strings.TrimRight(string(o), "\n")
	if err != nil {
		return true, fmt.Errorf(dir, ": git rev-parse to get branch name failed")
	}
	// find out if there is even a version of this branch on origin?

	// find out if there are any commits that local has that origin does not
	o, err = exec.Command("git", "rev-list", bn ,"--not", "origin/" + bn, "--count").Output()
	if err != nil {
		return true, fmt.Errorf(dir, ": git rev-list to compare to origin failed")
	}
	if strings.TrimRight(string(o), "\n") == "0" {
		return false, nil
	}
	return true, nil
}


func IsTreeClean(dir string) (bool, error) {
	os.Chdir(dir)
	o, err := exec.Command("git", "status").Output()
	status := string(o)
	if err != nil {
		return true, fmt.Errorf(dir, ": git status failed")
	}

	treeClean := regexp.MustCompile(`nothing to commit, working tree clean`)

	return treeClean.MatchString(status), nil
}

func IsRepo(dir string) (bool, error) {
	gitp := filepath.Join(dir, ".git")

	if info, err := os.Stat(gitp); os.IsNotExist(err) {
		return false, err
	} else if info.isDir() {
			return true, nil
	} else
		return false, nil
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

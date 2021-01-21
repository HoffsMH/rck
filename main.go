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

		repos, err := RepoList(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, r := range repos {
			rs, err := CheckRepo(r)
			if err != nil {
				return
			}

			if !rs {
				fmt.Println(r)
			}
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
		return false, err
	}
	remotes := len(strings.Split(string(o), "\n"))
	if remotes < 2 {
		return true, err
	}

	// fetch latest information from remotes
	_, err = exec.Command("git", "fetch").Output()

	if err != nil {
		return false, err
	}

	// gather git command outputs
	o, err = exec.Command("git", "status").Output()
	status := string(o)
	if err != nil {
		return false, err
	}

	// check that the repo is up to date
	upToDate, err := IsRepoUpToDate(status)
	if err != nil {
		return false, err
	}
	return upToDate, err
}

func IsRepoUpToDate(status string) (bool, error) {
	return IsBranchUpToDate(status) && IsTreeClean(status), nil
}

func IsBranchUpToDate(status string) bool {
	// confirm that there is "Your Branch is" before
	// asserting that the branch is up to date
	info := regexp.MustCompile(`Your branch is`)
	uptoDate := regexp.MustCompile(`Your branch is up to date with 'origin`)
	if !info.MatchString(status) {
		return true
	}
	return uptoDate.MatchString(status)
}

func IsTreeClean(status string) bool {
	treeClean := regexp.MustCompile(`nothing to commit, working tree clean`)
	return treeClean.MatchString(status)
}

func IsRepo(dir string) (bool, error) {
	gitp := filepath.Join(dir, ".git")

	if _, err := os.Stat(gitp); os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
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

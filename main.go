package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var rootCmd = &cobra.Command{
	Use:   "rck",
	Short: "repository check",
	Long: `check all repos recursively from a given directory
					for their up-to-dateness with origin`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := filepath.Abs(args[0])
		var repos []string

		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				fmt.Println(path, info.Size())
				return nil
			})

		if err != nil {
			return
		}

		o, err := CheckRepo(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(o)
	},
}

func CheckRepo(dir string) (bool, error) {
	if _, err := IsRepo(dir); err != nil {
		return false, err
	}

	_, err := exec.Command("git", "fetch").Output()
	upToDate, err := IsRepoUpToDate(dir)

	if err != nil {
		return false, err
	}
	return upToDate, nil
}

func IsRepoUpToDate(dir string) (bool, error) {
	os.Chdir(dir)
	output, err := exec.Command("git", "status").Output()
	status := string(output)

	if err != nil {
		return false, err
	}
	return IsBranchUpToDate(status) && IsTreeClean(status), nil
}

// Your branch is ahead of 'origin/package' by 1 commit.
// Your branch is up to date with 'origin/package'.
// Your branch and 'origin/master' have diverged,
func IsBranchUpToDate(status string) bool {
	branchUpToDate := regexp.MustCompile(`Your branch is up to date with 'origin`)
	return branchUpToDate.MatchString(status)
}

func IsTreeClean(status string) bool {
	treeClean := regexp.MustCompile(`nothing to commit, working tree clean`)
	return treeClean.MatchString(status)
}

func IsRepo(dir string) (string, error) {
	gitp := filepath.Join(dir, ".git")

	if _, err := os.Stat(gitp); os.IsNotExist(err) {
		return dir, err
	} else {
		return dir, nil
	}
}

var prependDateCmd = &cobra.Command{
	Use:     "prependdate",
	Aliases: []string{"pd"},
	Short:   "Downloads an episode of a podcast.",
	Long:    ``,
	Args:    cobra.MinimumNArgs(1),
	Run:     prependDate,
}

func prependDate(cmd *cobra.Command, args []string) {
	fmt.Println("prependDate")
	fmt.Println(args[0])
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

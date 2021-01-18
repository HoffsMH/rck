package main

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rck",
	Short: "repository check",
	Long: `check all repos recursively from a given directory for their up-to-dateness with origin`,
	Run: func(cmd *cobra.Command, args []string) {
		fp, err := filepath.Abs(args[0]);
		_ = err
		fmt.Println(fp)
	},
}

var prependDateCmd = &cobra.Command{
	Use:     "prependdate",
	Aliases: []string{"pd"},
	Short:   "Downloads an episode of a podcast.",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run:  prependDate,
}

func prependDate(cmd *cobra.Command, args []string) {
	fmt.Println("prependDate");
	fmt.Println(args[0]);
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
}

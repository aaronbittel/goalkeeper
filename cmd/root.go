package cmd

import (
	"os"

	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

var (
	tomlConfig  pkg.TomlDocument
	csvFilename string
	tasks       []*pkg.Task
	lastTask    *pkg.Task = nil
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A Cli tool for keeping track of progress.",
	Long: `To keep track for your programming journey progress.
	Set a goal and keep track your time spent on projects and programming languages.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	tomlConfig = pkg.LoadTomlConfig()
	csvFilename = tomlConfig.ConfigSection.Filename
	tasks = pkg.LoadTasks(csvFilename)

	if len(tasks) > 0 {
		lastTask = tasks[len(tasks)-1]
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// TODO: Do this

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
}

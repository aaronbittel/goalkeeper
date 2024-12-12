package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

var (
	tomlConfig pkg.TomlDocument
	tasks      []*pkg.Task
	lastTask   *pkg.Task
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goalkeeper",
	Short: "A Cli tool for keeping track of progress.",
	Long: `To keep track for your programming journey progress.
	Set a goal and keep track your time spent on projects and programming languages.`,
	PersistentPreRun: rootPreRun,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func rootPreRun(cmd *cobra.Command, args []string) {
	var err error
	tomlConfig, err = pkg.LoadTomlConfig()

	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("could not parse config.toml: %v", err)
		}

		// config file does not exist -> create config file
		pkg.CreateProjectDir(pkg.DefaultPath())
		err := pkg.CreateTomlFile(pkg.DefaultTomlConfig())
		if err != nil {
			log.Fatal(err)
		}
	}

	csvFilename := tomlConfig.ConfigSection.Filename

	tasks, err = pkg.LoadTasks(csvFilename)
	if err != nil {
		path := filepath.Join(pkg.DefaultPath() + csvFilename)
		_, err = os.Create(path)
		if err != nil {
			log.Fatalf("error creating csv file: %v", err)
		}
	}

	if len(tasks) > 0 {
		lastTask = tasks[len(tasks)-1]
	}
}

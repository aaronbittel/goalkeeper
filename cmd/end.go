package cmd

import (
	"fmt"

	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

// endCmd represents the end command
var endCmd = &cobra.Command{
	Use:   "end",
	Short: "Ends a running task.",
	Long: `Sets the end time for the currently running task and ends it.
	Now you can begin a new task with "start"`,
	Run:     runEnd,
	Aliases: []string{"stop"},
}

func runEnd(cmd *cobra.Command, args []string) {
	if len(tasks) == 0 || lastTask.IsFinished() {
		fmt.Println("First call 'start' to begin a new task")
		return
	}

	lastTask := tasks[len(tasks)-1]
	lastTask.Finish()
	pkg.SaveTasks(csvFilename, tasks)
}

func init() {
	rootCmd.AddCommand(endCmd)
}

package cmd

import (
	"fmt"
	"log"

	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

var (
	project  string
	language string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This starts a new task.",
	Long: `This starts a new task with for the given "Project" and "Language.
	The start time is set to now and the end time is TBD.
	Finish a task using the "end" command."`,
	Aliases: []string{"begin"},
	Run:     runStart,
}

func runStart(cmd *cobra.Command, args []string) {
	if len(tasks) > 0 && !lastTask.IsFinished() {
		fmt.Printf(
			"First call 'end' to finish the running task:\n\t %s (%s) started at: %s\n",
			lastTask.Project,
			lastTask.Language,
			lastTask.Start.Format("2006-01-02 15:04:05"),
		)
		return
	}

	task := pkg.NewTask(project, language)
	tasks = append(tasks, task)
	pkg.SaveTasks(csvFilename, tasks)

	log.Printf(
		"Successfully saved task %s (%s), started at: %s\n",
		task.Project,
		task.Language,
		task.Start.Format("2006-01-02 15:04:05"),
	)
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&project, "project", "p", "", "The name of the project of that task")
	startCmd.Flags().StringVarP(&language, "language", "l", "", "The programming language of that task")

	startCmd.MarkFlagRequired("project")
	startCmd.MarkFlagRequired("language")
}

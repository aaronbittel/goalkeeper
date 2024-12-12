package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	table "github.com/aaronbittel/goalkeeper/internal"
	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints the status of the current day.",
	Long: `This prints the status of the current day.
	This includes all the tasks (project, languages) as well as start and end times.
	This also sums up all the hours and shows the progress for today.`,
	Run: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("date", "d", "", "Retrieve the status of a specific day in the past")
	statusCmd.Flags().BoolP("yesterday", "y", false, "Retrieve yesterday's status")
	statusCmd.Flags().BoolP("percentage", "p", false, "Show the progress in percentage")
}

func runStatus(cmd *cobra.Command, args []string) {
	dateStr, err := cmd.Flags().GetString("date")
	if err != nil {
		log.Fatalf("[status] error getting date value: %v", err)
	}

	isYesterday, err := cmd.Flags().GetBool("yesterday")
	if err != nil {
		log.Fatalf("[status] error getting yesterday value: %v", err)
	}

	date := time.Now()
	if isYesterday {
		date = date.AddDate(0, 0, -1)
	}

	if !isYesterday && dateStr != "" {
		date, err = time.Parse("2.1.2006", dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"could not parse date: %s, please use format 'DD.MM.YYYY'\n",
				dateStr)
		}

		if date.After(time.Now()) {
			fmt.Fprintf(os.Stderr,
				"Time travel hasn't been invented yet\n%s is in the future\n",
				date.Format("2.1.2006"))
			return
		}
	}

	showPercentage, err := cmd.Flags().GetBool("percentage")
	if err != nil {
		log.Fatalf("[status] error getting percentage value: %v", err)
	}

	tasks := pkg.GetTasksForDate(tasks, date)
	if len(tasks) == 0 {
		fmt.Fprintf(os.Stderr, "There are no tasks for that day")
		return
	}

	printTasks(tasks, date, showPercentage)
}

func printTasks(tasks []*pkg.Task, date time.Time, showPercentage bool) {
	var totalDuration time.Duration

	tab := table.NewTable(
		table.NewHeader("Project").HeadingCentered(),
		table.NewHeader("Language", true),
		table.NewHeader("Start", true),
		table.NewHeader("End", true),
		table.NewHeader("Duration", true),
	).WithRoundedCorners().WithTitle(date.Format("Mon Jan 02 '06"))

	for _, t := range tasks {
		tab.AddRow([]string{
			t.Project, t.Language, t.Start.Format(pkg.TimeFormat),
			pkg.FormatTimeOrTBD(t.End, pkg.TimeFormat), formatDuration(t.Duration()),
		})
		totalDuration += t.Duration()
	}
	tab.AddSeperator()

	percentage := ""
	goalMinutes := tomlConfig.GoalsSection.Daily
	if showPercentage && goalMinutes != 0 {
		perc := totalDuration.Minutes() / float64(goalMinutes)
		percentage = fmt.Sprintf(" (%d%%)", int(perc*100))
	}

	tab.AddRow([]string{"", "", "", "", fmt.Sprintf(
		"%s%s",
		formatDuration(totalDuration),
		percentage)})

	fmt.Println(tab)
}

func formatDuration(dur time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(dur.Hours()), int(dur.Minutes())%60)
}

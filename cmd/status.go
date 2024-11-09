package cmd

import (
	"fmt"
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

func runStatus(cmd *cobra.Command, args []string) {
	tab := table.NewTable(
		table.NewHeader("Project").HeadingCentered(),
		table.NewHeader("Language", true),
		table.NewHeader("Start", true),
		table.NewHeader("End", true),
		table.NewHeader("Duration", true),
	).WithRoundedCorners().WithTitle(time.Now().Format("Mon Jan 02"))

	var totalDuration time.Duration

	todayTasks := pkg.GetTodayTasks(tasks)
	for _, t := range todayTasks {
		tab.AddRow([]string{
			t.Project, t.Language, t.Start.Format(pkg.TimeFormat),
			pkg.FormatTimeOrTBD(t.End, pkg.TimeFormat), formatDuration(t.Duration()),
		})
		totalDuration += t.Duration()
	}
	tab.AddSeperator()

	percentage := ""
	goalMinutes := tomlConfig.GoalsSection.Daily
	if goalMinutes != 0 {
		perc := totalDuration.Minutes() / float64(goalMinutes)
		percentage = fmt.Sprintf(" (%d%%)", int(perc*100))
	}

	tab.AddRow([]string{"", "", "", "", fmt.Sprintf(
		"%s%s",
		formatDuration(totalDuration),
		percentage)})

	fmt.Println(tab)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func formatDuration(dur time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(dur.Hours()), int(dur.Minutes())%60)
}

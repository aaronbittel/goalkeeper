package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

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
	var (
		todayTasks  = pkg.GetTodayTasks(tasks)
		amountToday time.Duration
		amount      time.Duration

		w = new(tabwriter.Writer)
	)

	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	_, err := fmt.Fprintf(w, "Project\tLanguage\tStart\tEnd\tDuration\n")
	if err != nil {
		log.Fatalf("error writing to tabwriter: %v", err)
	}

	for _, t := range todayTasks {
		amount = t.Duration()
		amountToday += amount

		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			t.Project, t.Language,
			pkg.FormatTimeOrTBD(t.Start, pkg.TimeFormat),
			pkg.FormatTimeOrTBD(t.End, pkg.TimeFormat),
			formatDuration(amount))
		if err != nil {
			log.Fatalf("error writing to tabwriter: %v", err)
		}
	}

	percentage := ""
	goalMinutes := tomlConfig.ConfigSection.GoalMinutes
	if goalMinutes != 0 {
		perc := amountToday.Minutes() / float64(goalMinutes)
		percentage = fmt.Sprintf(" (%2.2f %%)", perc*100)
	}

	_, err = fmt.Fprintf(w, "\t\t\t\t%s\n",
		fmt.Sprintf("%s%s", formatDuration(amountToday), percentage),
	)
	if err != nil {
		log.Fatalf("error writing to tabwriter: %v", err)
	}

	err = w.Flush()
	if err != nil {
		log.Fatalf("error flushing tabwriter: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func formatDuration(dur time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(dur.Hours()), int(dur.Minutes())%60)
}

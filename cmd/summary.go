package cmd

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"time"

	table "github.com/aaronbittel/goalkeeper/internal"
	"github.com/aaronbittel/goalkeeper/pkg"
	"github.com/spf13/cobra"
)

var (
	projectSummary  bool
	languageSummary bool
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Overview of this week's progress",
	// Long:  `Overview for this week's progress by day.`,
	Run: runSummary,
}

func runSummary(cmd *cobra.Command, args []string) {
	project, err := cmd.Flags().GetBool("project")
	if err != nil {
		log.Fatalf("[summary] error getting project value: %v", err)
	}

	language, err := cmd.Flags().GetBool("language")
	if err != nil {
		log.Fatalf("[summary] error getting project value: %v", err)
	}

	ascending, err := cmd.Flags().GetBool("ascending")
	if err != nil {
		log.Fatalf("[summary] error getting project value: %v", err)
	}

	if !project && !language {
		summaryWeek()
	}

	if project {
		summaryProjects(ascending)
	}

	if language {
		summaryLanguages(ascending)
	}

}

func init() {
	rootCmd.AddCommand(summaryCmd)

	summaryCmd.Flags().BoolP("project", "p", false, "Show project summary")
	summaryCmd.Flags().BoolP("language", "l", false, "Show language summary")
	summaryCmd.Flags().BoolP("ascending", "a", false, "Show output in ascending order")
}

func summaryWeek() {
	summary := make(map[time.Time][]*pkg.Task)
	for _, t := range slices.Backward(tasks) {
		if t.Start.Weekday() == time.Sunday {
			break
		}
		year, month, day := t.Start.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		summary[date] = append(summary[date], t)
	}

	printSummary(summary)
}

func printSummary(summary map[time.Time][]*pkg.Task) {
	weekdays := make([]time.Time, 0, len(summary))
	for t := range summary {
		weekdays = append(weekdays, t)
	}

	sort.Slice(weekdays, func(i, j int) bool {
		return weekdays[i].Before(weekdays[j])
	})

	table := table.NewTable(
		table.NewHeader("Date", false),
		table.NewHeader("Project", false),
		table.NewHeader("Language", true),
		table.NewHeader("Duration", false),
	)
	table.WithRoundedCorners()

	for j, weekday := range weekdays {
		// var amountToday time.Duration

		for i, t := range slices.Backward(summary[weekday]) {

			var weekdayStr string
			if i == len(summary[weekday])-1 {
				weekdayStr = weekday.Format(pkg.DateFormat)
			}

			table.AddRow([]string{
				weekdayStr,
				t.Project,
				t.Language,
				formatDuration(t.Duration()),
			})
		}

		if j != len(weekdays)-1 {
			table.AddRow([]string{"", "", "", ""})
		}
	}

	fmt.Println(table.String())
}

func summaryProjects(ascending bool) {
	projectTasks := map[string]time.Duration{}
	for _, t := range tasks {
		projectTasks[t.Project] += t.Duration()
	}

	projectNames := make([]string, 0, len(projectTasks))
	for k := range projectTasks {
		projectNames = append(projectNames, k)
	}

	if !ascending {
		sort.Slice(projectNames, func(i, j int) bool {
			return projectTasks[projectNames[i]] > projectTasks[projectNames[j]]
		})
	} else {
		sort.Slice(projectNames, func(i, j int) bool {
			return projectTasks[projectNames[i]] < projectTasks[projectNames[j]]
		})
	}

	tab := table.NewTable(
		table.NewHeader("Projects").HeadingCentered(),
		table.NewHeader("Duration", true),
	).WithRoundedCorners()

	for _, name := range projectNames {
		tab.AddRow([]string{name, formatDuration(projectTasks[name])})
	}

	fmt.Println(tab)
}

func summaryLanguages(ascending bool) {
	languageTasks := map[string]time.Duration{}
	for _, t := range tasks {
		languageTasks[t.Language] += t.Duration()
	}

	languageNames := make([]string, 0, len(languageTasks))
	for k := range languageTasks {
		languageNames = append(languageNames, k)
	}

	if !ascending {
		sort.Slice(languageNames, func(i, j int) bool {
			return languageTasks[languageNames[i]] > languageTasks[languageNames[j]]
		})
	} else {
		sort.Slice(languageNames, func(i, j int) bool {
			return languageTasks[languageNames[i]] < languageTasks[languageNames[j]]
		})
	}

	tab := table.NewTable(
		table.NewHeader("Languages", true),
		table.NewHeader("Duration", true),
	).WithRoundedCorners()

	for _, name := range languageNames {
		tab.AddRow([]string{name, formatDuration(languageTasks[name])})
	}

	fmt.Println(tab)
}

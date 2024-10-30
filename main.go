package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	var (
		project  string
		language string
	)

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmd.StringVar(&project, "project", "", "Sets the project for the task")
	startCmd.StringVar(&language, "language", "", "Sets the language for the task")

	flag.NewFlagSet("end", flag.ExitOnError)
	flag.NewFlagSet("status", flag.ExitOnError)

	if len(os.Args) < 2 {
		// TODO: Make this nicer
		log.Fatal("not enough arguments")
	}

	tomlConfig := loadTomlConfig()
	csvFilename := tomlConfig.ConfigSection.Filename
	tasks := loadTasks(csvFilename)
	var lastTask *Task = nil
	if len(tasks) > 0 {
		lastTask = tasks[len(tasks)-1]
	}

	switch os.Args[1] {
	case "start":
		if len(tasks) > 0 && !lastTask.IsFinished() {
			fmt.Printf(
				"First call 'end' to finish the running task:\n\t %s (%s) started at: %s\n",
				lastTask.Project,
				lastTask.Language,
				lastTask.Start.Format("2006-01-02 15:04:05"),
			)
			return
		}

		// TODO: Make this if its not given it is asked for
		startCmd.Parse(os.Args[2:])
		if project == "" {
			fmt.Fprintf(os.Stderr, "project name is required\n")
			return
		}
		if language == "" {
			fmt.Fprintf(os.Stderr, "language is required\n")
			return
		}

		task := NewTask(project, language)
		tasks = append(tasks, task)
		saveTasks(csvFilename, tasks)

		log.Printf(
			"Successfully saved task %s (%s), started at: %s\n",
			task.Project,
			task.Language,
			task.Start.Format("2006-01-02 15:04:05"),
		)

	case "end":
		if len(tasks) == 0 || lastTask.IsFinished() {
			fmt.Println("First call 'start' to begin a new task")
			return
		}

		lastTask := tasks[len(tasks)-1]
		lastTask.Finish()
		saveTasks(csvFilename, tasks)
	case "status":
		var (
			todayTasks  = getTodayTasks(tasks)
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
			if t.End.IsZero() {
				amount = time.Now().Sub(t.Start)
			} else {
				amount = t.End.Sub(t.Start)
			}
			amountToday += amount

			_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				t.Project, t.Language,
				formatTimeOrTBD(t.Start, TimeFormat),
				formatTimeOrTBD(t.End, TimeFormat),
				formatAmount(amount))
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
			fmt.Sprintf("%s%s", formatAmount(amountToday), percentage),
		)
		if err != nil {
			log.Fatalf("error writing to tabwriter: %v", err)
		}

		err = w.Flush()
		if err != nil {
			log.Fatalf("error flushing tabwriter: %v", err)
		}

	default:
		fmt.Printf("[ERROR] unknown subcommand '%s'\n", os.Args[1])
		os.Exit(1)
	}

}

func formatAmount(dur time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(dur.Hours()), int(dur.Minutes())%60)
}

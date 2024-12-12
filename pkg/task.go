package pkg

import (
	"fmt"
	"log"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04"
	DateFormat     = "2006-01-02"
)

type Task struct {
	Project  string
	Language string
	Start    time.Time
	End      time.Time
}

func NewTask(project, language string) *Task {
	return &Task{
		Project:  project,
		Language: language,
		Start:    time.Now(),
	}
}

func (t Task) Fields() []string {
	return []string{t.Project, t.Language, t.Start.Format("2006-01-02 15:04:05"),
		FormatTimeOrTBD(t.End, DateTimeFormat)}
}

func (t Task) String() string {
	return fmt.Sprintf(
		"Project: %q, Language: %q: Started: %s, Ended: %s",
		t.Project,
		t.Language,
		t.Start.Format("2006-01-02 15:04:05"),
		FormatTimeOrTBD(t.End, DateTimeFormat),
	)
}

func FormatTimeOrTBD(t time.Time, format string) string {
	if t.IsZero() {
		return "TBD"
	}
	return t.Format(format)
}

func (t Task) IsFinished() bool {
	return !t.End.IsZero()
}

func FromFields(fields []string) *Task {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatalf("error loading time zone location: %v", err)
	}

	start, err := time.ParseInLocation("2006-01-02 15:04:05", fields[2], location)
	if err != nil {
		log.Fatalf("[start time] Error parsing time (%s): %v", fields[2], err)
	}

	var end time.Time

	if fields[3] == "TBD" {
		end = time.Time{}
	} else {
		end, err = time.ParseInLocation("2006-01-02 15:04:05", fields[3], location)
		if err != nil {
			log.Fatalf("[end time] Error parsing time (%s): %v", fields[3], err)
		}
	}

	return &Task{
		Project:  fields[0],
		Language: fields[1],
		Start:    start,
		End:      end,
	}
}

func (t *Task) Finish() {
	t.End = time.Now()
}

func GetTasksForDate(tasks []*Task, t time.Time) []*Task {
	todayTasks := []*Task{}
	for _, task := range tasks {
		if SameDay(t, task.Start) {
			todayTasks = append(todayTasks, task)
		}
	}
	return todayTasks
}

func (t Task) Duration() time.Duration {
	endtime := t.End
	if t.End.IsZero() {
		endtime = time.Now()
	}

	return endtime.Sub(t.Start)
}

func SameDay(ref, t time.Time) bool {
	nYear, nMonth, nDay := ref.Date()
	year, month, day := t.Date()
	return nYear == year && nMonth == month && nDay == day
}

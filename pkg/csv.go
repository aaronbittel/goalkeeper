package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func LoadTasks(filename string) []*Task {
	var f *os.File
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(filename)
			if err != nil {
				log.Fatalf("error creating csv file: %v", err)
			}
		} else {
			log.Fatalf("[load] error opening file %s: %v", filename, err)
		}
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalf("error reading csv: %v", err)
	}

	if len(records) == 0 {
		return []*Task{}
	}

	// Dismiss column names
	records = records[1:]

	tasks := make([]*Task, 0, len(records))

	for _, record := range records {
		tasks = append(tasks, FromFields(record))
	}

	return tasks
}

// TODO: Save to new file / backup old file, if error occurs restore old file
func SaveTasks(filename string, tasks []*Task) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("[save] error opening file %s: %v", filename, err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)

	if err := writer.Write([]string{"Category", "Title", "Start", "End"}); err != nil {
		log.Fatalf("error writing csv column names to file")
	}

	for _, task := range tasks {
		if writer.Write(task.Fields()); err != nil {
			log.Fatalf("error writing %s to %s: %v", task.Project, filename, err)
		}
	}
	writer.Flush()

	if writer.Error() != nil {
		log.Fatalf("error flushing csv writer: %v", err)
	}

	fmt.Printf("Successfully saved record to %s\n", filename)
}

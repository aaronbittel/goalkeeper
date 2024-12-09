package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func LoadTasks(filename string) ([]*Task, error) {
	path := filepath.Join(DefaultPath(), filename)
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		log.Fatalf("[load] error opening file %s: %v", filename, err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		// log.Fatalf("error reading csv: %v", err)
		return nil, err
	}

	if len(records) == 0 {
		return []*Task{}, nil
	}

	// Dismiss column names
	records = records[1:]

	tasks := make([]*Task, 0, len(records))

	for _, record := range records {
		tasks = append(tasks, FromFields(record))
	}

	return tasks, nil
}

// TODO: Save to new file / backup old file, if error occurs restore old file
func SaveTasks(filename string, tasks []*Task) {
	path := filepath.Join(DefaultPath(), filename)
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
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

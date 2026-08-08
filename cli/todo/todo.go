package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item struct reps a ToDo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List reps a list to ToDo item
type List []item

// Add creates a new ToDo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete method marks the ToDo item as completed
// by setting Done = true and CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l

	ls.isExist(i)

	// Adjusting index for zero based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete method removes an item from the ToDo list
func (l *List) Delete(i int) error {
	ls := *l

	ls.isExist(i)

	// Adjust index for zero based index
	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save method encodes the list as JSON and saves it
// using the provided file name
func (l *List) Save(fname string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(fname, js, 0644)
}

// Get method opens the provided file name, decodes
// the JSON and parses it into a List
func (l *List) Get(fname string) error {
	file, err := os.ReadFile(fname)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

func (l *List) String() string {
	formatted := ""

	for i, t := range *l {
		prefix := "  "
		if t.Done {
			prefix = "X "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, i+1, t.Task)
	}

	return formatted
}

// Utility function to check for i not found in the list
func (l *List) isExist(i int) error {
	if i <= 0 || i >= len(*l) {
		return fmt.Errorf("Item %d, does not exit", i)
	}

	return nil
}

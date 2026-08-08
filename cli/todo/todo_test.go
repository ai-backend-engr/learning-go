package todo_test

import (
	"example/todo"
	"os"
	"testing"
)

// TestAdd tests the Add method of the List type
func TestAdd(t *testing.T) {
	l := todo.List{}

	task := "New todo"
	l.Add(task)

	if task != l[0].Task {
		t.Errorf("Expected %q, got %q instead", task, l[0].Task)
	}
}

// TestComplete tests the Complete method of the List type
func TestComplete(t *testing.T) {
	l := todo.List{}

	task := "New todo"
	l.Add(task)

	if task != l[0].Task {
		t.Errorf("Expected %q, got %q instead", task, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New task should not be completed")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("New task should be completed")
	}
}

// TestDelete test the Delete method of the List type
func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"New task 1",
		"New task 2",
		"New task 3",
	}

	for _, v := range tasks {
		l.Add(v)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Expected %q, got %q instead", tasks[0], l[0].Task)
	}

	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead", tasks[2], l[2].Task)
	}
}

// TestSaveGet tests the Save and Get method of the List type
func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	task := "New task"
	l1.Add(task)

	if l1[0].Task != task {
		t.Errorf("Expected %q, got %q instead", task, l1[0].Task)
	}

	tempF, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempF.Name())

	if err := l1.Save(tempF.Name()); err != nil {
		t.Errorf("Error saving list to file: %s", err)
	}

	if err := l2.Get(tempF.Name()); err != nil {
		t.Errorf("Error getting list from file: %s", err)
	}

	if l2[0].Task != l1[0].Task {
		t.Errorf("Task %q should match %q task", l1[0].Task, l2[0].Task)
	}
}

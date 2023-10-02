package main

// Task Item Interface

type status int

const (
	todo status = iota
	inProgress
	done
)

// Task data model
type Task struct {
	status      status
	title       string
	description string
}

// FilterValue in task
func (t Task) FilterValue() string {
	return t.title
}

// Title of task
func (t Task) Title() string {
	return t.title
}

// Description of task
func (t Task) Description() string {
	return t.description
}

// Next Task
func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

func NewTask(status status, title, description string) Task {
	return Task{title: title, description: description, status: status}
}

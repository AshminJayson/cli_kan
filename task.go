package main

import (
	"encoding/json"
	"fmt"
)

// Task Item Interface

type status int

const (
	todo status = iota
	inProgress
	done
)

// Task data model
type Task struct {
	Status          status
	TaskTitle       string
	TaskDescription string
}

// FilterValue in task
func (t Task) FilterValue() string {
	return t.TaskTitle
}

// Title of task
func (t Task) Title() string {
	return t.TaskTitle
}

// Description of task
func (t Task) Description() string {
	return t.TaskDescription
}

// Next Task
func (t *Task) Next() {
	if t.Status == done {
		t.Status = todo
	} else {
		t.Status++
	}
}

func NewTask(status status, title, description string) Task {
	return Task{TaskTitle: title, TaskDescription: description, Status: status}
}

func (t Task) getJson() []byte {
	content, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))
	return content
}

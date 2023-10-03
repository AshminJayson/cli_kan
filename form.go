package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	create int = iota
	edit
)

// Form Model
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
	operation   int
}

func NewForm(focused status, listitem list.Item) *Form {
	form := &Form{focused: focused, operation: create}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()

	if listitem != nil {
		form.operation = edit
		task := listitem.(Task)
		form.title.SetValue(task.TaskTitle)
		form.description.SetValue(task.TaskDescription)
	}
	return form
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) CreateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value())
	return task
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "enter":
				if m.title.Focused() {
					m.title.Blur()
					m.description.Focus()
					return m, textarea.Blink
				}

				models[form] = m
				return models[mainModel], m.CreateTask
			case "esc":
				if m.description.Focused() {
					m.description.Blur()
					m.title.Focus()
					return m, textinput.Blink
				}

				if m.operation == edit {
					return models[mainModel], m.CreateTask
				}
				return models[mainModel], nil
			}
		}
	}

	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	}

	m.description, cmd = m.description.Update(msg)
	return m, cmd

}

func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}

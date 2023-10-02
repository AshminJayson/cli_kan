package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// Main Model

const divisor = 4

// Model for tea
type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

// Model Manager
var models []tea.Model

const (
	mainModel status = iota
	form
)

// Styling

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// Next list focus
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

// Prev list focus
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-10)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	m.lists[todo].Title = " To Do"
	m.lists[todo].SetItems(([]list.Item{
		Task{status: todo, title: "buy milk", description: "straw berry milk"},
		Task{status: todo, title: "buy milk", description: "banana milk"},
		Task{status: todo, title: "buy milk", description: "chocolate milk"},
	}))
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems(([]list.Item{
		Task{status: inProgress, title: "buy phone", description: "straw berry milk"},
		Task{status: inProgress, title: "buy laptop", description: "banana milk"},
		Task{status: inProgress, title: "buy bus", description: "chocolate milk"},
	}))
	m.lists[done].Title = " Done"
	m.lists[done].SetItems(([]list.Item{
		Task{status: done, title: "buy cat", description: "straw berry cat"},
		Task{status: done, title: "buy dog", description: "banana dog"},
		Task{status: done, title: "buy cow", description: "chocolate cow"},
	}))
}

// NewMainModel func for tea
func NewMainModel() *Model {
	return &Model{err: nil}
}

// MoveToNext task
func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()

	if selectedItem == nil {
		return nil
	}
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

// DeleteTask that is currently selected
func (m *Model) DeleteTask() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return nil
	}

	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	return nil

}

// Init func for tea
func (m Model) Init() tea.Cmd {
	return nil
}

// Update func for tea
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			// columnStyle.Width(msg.Width / divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		case "delete":
			return m, m.DeleteTask
		case "n":
			models[mainModel] = m
			models[form] = NewForm(m.focused)
			return models[form].Update(nil)

		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

// View func for tea
func (m Model) View() string {

	if m.quitting {
		return ""
	}

	if !m.loaded {
		return "loading.."
	}

	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()

	switch m.focused {
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case inProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneView),
		)

	}

}

// Form Model
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
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
				} else {
					models[form] = m
					return models[mainModel], m.CreateTask
				}
			}
		}
	}

	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}

func main() {
	models = []tea.Model{NewMainModel(), NewForm(todo)}
	m := models[mainModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

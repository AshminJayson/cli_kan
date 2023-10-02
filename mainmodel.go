package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

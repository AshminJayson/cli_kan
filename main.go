package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
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

// Main Model

const divisor = 4

// Model for tea
type Model struct {
	focused status
	lists   []list.Model
	err     error
	loaded  bool
}

// Styling

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

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

// New func for tea
func New() *Model {
	return &Model{err: nil}
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
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

// View func for tea
func (m Model) View() string {
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

	}

}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

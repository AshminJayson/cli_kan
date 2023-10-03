package main

import (
	"encoding/json"
	"log"
	"os"

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
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("6"))
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

func CreateTaskFile() {
	Home, _ := os.UserHomeDir()
	f, err := os.Create(Home + "\\tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("[]"))
	f.Close()
}

// GetTasksFromFile for getting tasks
func GetTasksFromFile() []Task {
	Home, _ := os.UserHomeDir()
	_, err := os.Stat(Home + "\\tasks.json")
	if err != nil {
		CreateTaskFile()
	}
	content, err := os.ReadFile(Home + "\\tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	tasks := []Task{}

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	return tasks
}

// WriteTasksToFile for writing to file
func (m Model) WriteTasksToFile() {
	CreateTaskFile()

	tasks := []Task{}
	for _, status := range [3]status{todo, inProgress, done} {
		for _, listItem := range m.lists[status].Items() {
			task := listItem.(Task)
			tasks = append(tasks, task)
		}
	}
	content := getJson(tasks)
	Home, _ := os.UserHomeDir()

	os.WriteFile(Home+"\\tasks.json", content, 0644)
}

func doneListDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.NormalDesc.Strikethrough(true)
	d.Styles.NormalDesc.Foreground(lipgloss.Color("8"))
	d.Styles.SelectedDesc.Strikethrough(true)
	d.Styles.SelectedDesc.Foreground(lipgloss.Color("8"))
	d.Styles.NormalTitle.Strikethrough(true)
	d.Styles.NormalTitle.Foreground(lipgloss.Color("8"))
	d.Styles.SelectedTitle.Strikethrough(true)
	d.Styles.SelectedTitle.Foreground(lipgloss.Color("8"))

	return d
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-10)
	defaultList.Styles.Title.Foreground(lipgloss.Color("0"))
	defaultList.Styles.Title.Background(lipgloss.Color("10"))

	doneList := list.New([]list.Item{}, doneListDelegate(), width/divisor, height-10)
	defaultList.SetShowHelp(false)
	doneList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, doneList}

	tasks := GetTasksFromFile()
	m.lists[todo].Title = "To Do"
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].Title = " Done"

	for _, task := range tasks {
		switch task.Status {
		case todo:
			m.lists[todo].InsertItem(len(m.lists[todo].Items()), task)
		case inProgress:
			m.lists[inProgress].InsertItem(len(m.lists[inProgress].Items()), task)
		case done:
			m.lists[done].InsertItem(len(m.lists[done].Items()), task)
		}
	}
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
	m.lists[selectedTask.Status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.Status].InsertItem(len(m.lists[selectedTask.Status].Items())-1, list.Item(selectedTask))
	return m
}

func (m *Model) MoveToPrevious() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return nil
	}

	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.Status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Previous()
	m.lists[selectedTask.Status].InsertItem(len(m.lists[selectedTask.Status].Items())-1, list.Item(selectedTask))
	return m

}

// DeleteTask that is currently selected
func (m *Model) DeleteTask() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return nil
	}

	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.Status].RemoveItem(m.lists[m.focused].Index())
	return m

}

// Init func for tea
func (m Model) Init() tea.Cmd {
	return nil
}

// Update func for tea
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.loaded {
		m.WriteTasksToFile()
	}
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
		case "backspace":
			return m, m.MoveToPrevious
		case "delete":
			return m, m.DeleteTask
		case "n":
			models[mainModel] = m
			models[form] = NewForm(m.focused, nil)
			return models[form].Update(nil)
		case "e":
			models[mainModel] = m
			models[form] = NewForm(m.focused, m.lists[m.focused].SelectedItem())
			m.DeleteTask()
			return models[form].Update(tea.KeyDelete)
		}
	case Task:
		task := msg
		return m, m.lists[task.Status].InsertItem(len(m.lists[task.Status].Items()), task)
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

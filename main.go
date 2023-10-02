package main

import (
	"log"

	"encoding/json"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// FileWriter for writing to file
func FileWriter(content []byte) {
	os.WriteFile("tasks.json", content, 0644)
}

// FileReader for getting tasks
func FileReader() []Task {
	content, err := os.ReadFile("tasks.json")
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

func main() {
	models = []tea.Model{NewMainModel(), NewForm(todo)}
	m := models[mainModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

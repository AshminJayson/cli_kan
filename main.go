package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	models = []tea.Model{NewMainModel(), NewForm(todo)}
	m := models[mainModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

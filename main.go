package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	models = []tea.Model{NewMainModel(), NewForm(todo, nil)}
	m := models[mainModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

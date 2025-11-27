package handlers

import (
	"html/template"

	"tse-p3/simulation"
)

var (
	tmpl *template.Template
	MainSimulation	*simulation.Simulation
)


func Initialize(sim *simulation.Simulation) {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	MainSimulation = sim

}
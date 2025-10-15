// server.go
package website

import (
    "fmt"
    "html/template"
    "net/http"
    "tse-p2/simulation"
)

var (
    homeTemplate *template.Template
    Address      string
    Sim          *simulation.Simulation
    Hub          *hub
)

func Initialize(addr string, sim *simulation.Simulation) {
    Address = addr
    Sim = sim

    homeTemplate = template.Must(template.ParseFiles("templates/index.html"))

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        websocketHandler(w, r, sim)
    })
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    Hub = newHub()
    go Hub.Run()

}

func Begin() {
    if err := http.ListenAndServe(Address, nil); err != nil {
        fmt.Printf("Web server error: %v\n", err)
    }
}
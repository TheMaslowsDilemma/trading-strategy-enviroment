package website

import (
	"fmt"
	"html/template"
	"net/http"

	"tse-p3/website/handlers"
	"tse-p3/simulation"
)


func Initialize(sim *simulation.Simulation) {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	handlers.Initialize(sim)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var (
			method string
		)

		method = r.Method

		switch method {
		case http.MethodGet:
			handlers.LoginGET(w, r)
		case http.MethodPost:
			handlers.LoginPOST(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var (
			method string
		)

		method = r.Method

		switch method {
		case http.MethodGet:
			handlers.RegisterGET(w, r)
		case http.MethodPost:
			handlers.RegisterPOST(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", handlers.Logout)

	// === Protected routes (require login) ===
	http.HandleFunc("/", authMiddleware(handlers.ChartsHandler))
	http.HandleFunc("/ws", authMiddleware(handlers.WebsocketHandler))
}

func Begin(addr string) {
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Web server error: %v\n", err)
	}
}
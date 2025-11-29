package handlers

import (
	"fmt"
	"context"
	"net/http"

	"tse-p3/users"
	"tse-p3/website/sessions"
)

type formData struct {
	Title	string
	Error   string
	Success string
}


func RegisterGET(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", formData{ Title: "Register" })
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		username string
		password string
		ctx      context.Context
	)

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get form values
	username = r.FormValue("username")
	password = r.FormValue("password")

	if username == "" || password == "" {
		tmpl.ExecuteTemplate(w, "login.html", formData{
			Title: "Register",
			Error: "All fields are required",
		})
		return
	}

	ctx = context.Background()
	err = users.CreateUser(ctx, username, password, MainSimulation)
	fmt.Println(err)
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", formData{
			Title: "Register",
			Error: "Username already taken or server error",
		})
		return
	}

	http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", formData{ Title: "Login" } )
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	ctx := context.Background()
	user, err := users.GetUserByName(ctx, username)
	if err != nil || !user.ComparePassword(password) {
		tmpl.ExecuteTemplate(w, "login.html", formData{Title: "Login", Error: "Invalid username or password"})
		return
	}

	if err := sessions.Set(w, user.ID); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessions.Clear(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


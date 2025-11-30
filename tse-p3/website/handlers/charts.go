package handlers

import (
	"fmt"
	"context"
	"net/http"
	"tse-p3/users"
)

func ChartsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx		context.Context
		user	users.User
		userOK	bool
		data	map[string]interface{}
		value	interface{}
		err		error
	)
	
	ctx = r.Context()
	
	// --- Get username from context ---
	value = ctx.Value("user")
	if value != nil {
		user, userOK = value.(users.User)
	}
	if !userOK {
		fmt.Println("No user found on charts page")
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	data = map[string]interface{}{
		"Username": user.Name,
	}


	err = tmpl.ExecuteTemplate(w, "charts.html", data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}
}
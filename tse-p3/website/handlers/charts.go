package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"tse-p3/users"
)

var (
	homeTmpl *template.Template
)

func init() {
	var err error
	homeTmpl, err = template.ParseFiles("templates/charts.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse home template: %v", err))
	}
}

func ChartsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx				context.Context
		username		string
		usernameOK		bool
		subscriptions	[]users.DataSubscription // ← fixed type name
		subscriptionsOK	bool
		data			map[string]interface{}
		value			interface{}
		err				error
	)

	ctx = r.Context()

	// --- Get username from context ---
	value = ctx.Value("user.name")
	if value != nil {
		username, usernameOK = value.(string)
	}
	if !usernameOK || username == "" {
		username = "unknown"
	}


	value = ctx.Value("user.subscriptions")
	if value != nil {
		subscriptions, subscriptionsOK = value.([]users.DataSubscription)
	}

	if !subscriptionsOK {
		subscriptions = []users.DataSubscription{}
	}


	data = map[string]interface{}{
		"Username":      username,
		"Subscriptions": subscriptions,
	}


	err = homeTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}
}
package website

import (
	"context"
	"html/template"
	"net/http"

	"tse-p3/users"
	"tse-p3/website/sessions"
)

var (
	templates *template.Template
	Address   string
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			usr_id	int64
			ok		bool
			usr		users.User
			ctx		context.Context
			err		error
		)

		usr_id, ok = sessions.Get(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx = r.Context()

		usr, err = users.GetUserById(ctx, usr_id)
		
		if err != nil {
			sessions.Clear(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx = context.WithValue(ctx, "user", usr)

		next(w, r.WithContext(ctx))
	}
}
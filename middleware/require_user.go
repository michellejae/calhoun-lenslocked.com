package middleware

import (
	"fmt"
	"net/http"

	"gitlab.com/michellejae/lenslocked.com/context"
	"gitlab.com/michellejae/lenslocked.com/models"
)

type RequireUser struct {
	models.UserService
}

// just have to create this cause some of our routes on main.go are handleFunc vs handlerFunc
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// use a closure to check cookies are there before going to specific routes
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		fmt.Println("User Found:", user)
		next(w, r)

	})
}

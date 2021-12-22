package middleware

import (
	"net/http"
	"strings"

	"gitlab.com/michellejae/lenslocked.com/context"
	"gitlab.com/michellejae/lenslocked.com/models"
)

type User struct {
	models.UserService
}

// just have to create this cause some of our routes on main.go are handleFunc vs handlerFunc
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// use a closure to check cookies are there before going to specific routes
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)

	})
}

// assumes that User middleware has already been run otherwise it will error
type RequireUser struct {
	User
}

// Apply assumes that User middlewere has already bee run
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// Apply assumes that User middlewere has already bee run

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return mw.User.ApplyFn(func(w http.ResponseWriter, r *http.Request) {
		// if the user is requesting a static asset or image we will not need to look up the current user
		// so we skip doing it
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})

}

// func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
// 	ourHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := context.User(r.Context())
// 		if user == nil {
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}
// 		next(w, r)
// 	})
// 	return mw.User.Apply(ourHandler)
// }

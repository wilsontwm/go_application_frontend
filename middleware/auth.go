package middleware

import (
	"app_frontend/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

var Logging = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//fmt.Println(time.Now(), ":", r.URL.Path, "@", r.Method)

			handler.ServeHTTP(w, r)
		})
	}
}

// Check if the user is authenticated
var CheckAuth = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCookie := controllers.ReadEncodedCookieHandler(w, r, "auth")

			if authCookie == "" {
				// Redirect to login page when the auth cookie is not set, at the same time same the current path in the cookie
				controllers.SetCookieHandler(w, r, "nextURL", r.URL.Path)

				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}

// Check if the user is logged in
var IsLoggedIn = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// list of the urls that do not require checking if it's logged in
			toSkip := []string{"/", "/noaccess", "/logout"}
			requestPath := r.URL.Path // current request path

			// Check if the request need authentication,
			// If not, then serve the request
			for _, value := range toSkip {
				if value == requestPath {
					handler.ServeHTTP(w, r)
					return
				}
			}

			authCookie := controllers.ReadEncodedCookieHandler(w, r, "auth")

			// If cookie has been set, then redirect to dashboard
			if authCookie != "" {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
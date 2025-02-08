package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html/"))))

	mux.Route("/web", func(mux chi.Router) {
		mux.Post("/auth", app.authenticate)
		mux.Get("/refresh-token", app.refreshUsingCookie)
		mux.Get("/logout", app.deleteRefreshCookie)

	})
	// authentication routes - auth handler, refresh handler
	mux.Post("/v1/auth", app.authenticate)
	mux.Post("/v1/refresh-token", app.refresh)

	// protected routes
	mux.Route("/v1/users", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/", app.allUsers)
		mux.Get("/{id}", app.getUser)
		mux.Delete("/{id}", app.deleteUser)
		mux.Put("/{id}", app.insertUser)
		mux.Patch("/", app.updateUser)
	})

	return mux
}

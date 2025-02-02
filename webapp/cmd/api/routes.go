package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.use(app.enableCORS)

	// authentication routes - auth handler, refresh handler
	mux.Post("/v1/auth", app.authenticate)
	mux.Post("/v1/refresh-token", app.refresh)

	// test handler
	mux.Get("/v1/test", func(w http.ResponseWriter, r *http.Request) {
		var payload = struct {
			Message string `json:"message"`
		}{
			Message: "Hello, world!",
		}
		_ = app.writeJSON(w, http.StatusOK, payload)
	})

	// protected routes
	mux.Route("/v1/users", func(mux chi.Router) {
		//mux.Use(app.authenticateUser)

		mux.Get("/", app.allUsers)
		mux.Get("/{id}", app.getUser)
		mux.Delete("/{id}", app.deleteUser)
		mux.Put("/{id}", app.insertUser)
		mux.Patch("/", app.updateUser)
	})

	return mux
}

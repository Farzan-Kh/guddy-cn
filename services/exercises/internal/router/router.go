package router

import (
	service "github.com/Farzan-kh/guddy-cn/exercises/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures the router for the exercises service
func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/exercises", service.GetExercises)
		r.Get("/program/{uuid}", service.GetProgram)
		r.Get("/completeProgram/{uuid}", service.GetCompleteProgram)
		r.Post("/program", service.PostProgram)
	})

	return r
}

package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ramnath.1998/apica-app/controllers"
)

func RunRoutes() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/cache", controllers.GetCache)
		r.Get("/cache/node", controllers.UpdateCache)
		r.Post("/cache/node", controllers.UpdateCache)

	})

	http.ListenAndServe(":8000", r)
}

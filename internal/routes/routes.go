package routes

import (
	"net/http"

	"article/internal/handler"
	"article/internal/response"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// InitRoutes initialises routes
func InitRoutes(app *handler.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// handling 404 page not found error
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.New().NotFound(w, http.StatusText(http.StatusNotFound))
	})

	// handling 405 method not allowed error
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.New().NotAllowed(w, http.StatusText(http.StatusMethodNotAllowed))
	})

	// route to handle article request
	r.Post("/articles", app.CreateArticle())
	r.Get("/articles/{article_id}", app.GetArticle())
	r.Get("/articles", app.GetArticles())

	return r
}

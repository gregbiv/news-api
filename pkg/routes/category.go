package routes

import (
	"github.com/go-chi/chi"
	"github.com/gregbiv/news-api/pkg/api"
	"github.com/gregbiv/news-api/pkg/api/category"
	"github.com/gregbiv/news-api/pkg/middleware"
	storageCategory "github.com/gregbiv/news-api/pkg/storage/category"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// RouteCategory registers category routes
func RouteCategory(urlExtractor api.URLExtractor, db *sqlx.DB) func(r chi.Router) {
	getter := storageCategory.NewGetter(db)
	storer := storageCategory.NewStorer(db)
	updater := storageCategory.NewUpdater(db)
	discarder := storageCategory.NewDiscarder(db)

	return func(r chi.Router) {
		r.With(
			middleware.JSONRequestSchema("create_category.json"),
			middleware.JSONDebugResponseSchema(map[int]string{
				http.StatusOK:         "create_category.json",
				http.StatusBadRequest: "error.json",
			}),
		).Post("/", category.NewPostCategoryHandler(getter, storer).ServeHTTP)
		r.Route("/{category_id}", func(r chi.Router) {
			r.With(
				middleware.JSONDebugResponseSchema(map[int]string{
					http.StatusOK:         "get_category.json",
					http.StatusNotFound:   "error.json",
					http.StatusBadRequest: "error.json",
				}),
			).Get("/", category.NewGetCategoryHandler(getter, urlExtractor).ServeHTTP)
			r.With(
				middleware.JSONRequestSchema("update_category.json"),
				middleware.JSONDebugResponseSchema(map[int]string{
					http.StatusOK:         "update_category.json",
					http.StatusBadRequest: "error.json",
				}),
			).Put("/", category.NewCategoryUpdateHandler(getter, updater, urlExtractor).ServeHTTP)
			r.With(
				middleware.JSONDebugResponseSchema(map[int]string{
					http.StatusNotFound:   "error.json",
					http.StatusBadRequest: "error.json",
				}),
			).Delete("/", category.NewDiscardCategoryHandler(getter, discarder, urlExtractor).ServeHTTP)
		})
	}
}

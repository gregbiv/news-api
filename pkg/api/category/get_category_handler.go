package category

import (
	"github.com/go-chi/render"
	"github.com/gregbiv/news-api/pkg/api"
	storageCategory "github.com/gregbiv/news-api/pkg/storage/category"
	"net/http"
)

type (
	getCategoryHandler struct {
		getter       storageCategory.Getter
		urlExtractor api.URLExtractor
	}
)

// NewGetCategoryHandler init and returns an instance of getCategoryHandler
func NewGetCategoryHandler(
	categoryGetter storageCategory.Getter,
	urlExtractor api.URLExtractor,
) http.Handler {
	return &getCategoryHandler{
		getter:       categoryGetter,
		urlExtractor: urlExtractor,
	}
}

func (h *getCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate query params
	itemID, err := h.urlExtractor.UUIDFromRoute(r, "category_id")
	if err != nil {
		api.NotFound(w, r)
		return
	}

	dbCategory, err := h.getter.GetCategoryByID(itemID.String())
	if err != nil {
		if err == storageCategory.ErrCategoryNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	modelCategory := category{}
	err = modelCategory.fromDB(dbCategory)
	if err != nil {
		if err == storageCategory.ErrCategoryNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, modelCategory)
}

package category

import (
	"github.com/go-chi/render"
	"github.com/gregbiv/news-api/pkg/api"
	storageCategory "github.com/gregbiv/news-api/pkg/storage/category"
	"net/http"
)

type discardCategoryHandler struct {
	getter       storageCategory.Getter
	discarder    storageCategory.Discarder
	urlExtractor api.URLExtractor
}

// NewDiscardCategoryHandler init and returns an instance of discardCategoryHandler
func NewDiscardCategoryHandler(
	getter storageCategory.Getter,
	discarder storageCategory.Discarder,
	urlExtractor api.URLExtractor,
) http.Handler {
	return &discardCategoryHandler{
		getter:       getter,
		discarder:    discarder,
		urlExtractor: urlExtractor,
	}
}

func (h *discardCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	err = h.discarder.Discard(dbCategory.CategoryID)
	if err != nil {
		if err == storageCategory.ErrCategoryNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

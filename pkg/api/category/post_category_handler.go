package category

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/gregbiv/news-api/pkg/api"
	storageCategory "github.com/gregbiv/news-api/pkg/storage/category"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type postCategoryHandler struct {
	getter storageCategory.Getter
	storer storageCategory.Storer
}

// NewPostCategoryHandler init and returns an instance of postCategoryHandler
func NewPostCategoryHandler(
	getter storageCategory.Getter,
	storer storageCategory.Storer,
) http.Handler {
	return &postCategoryHandler{
		getter: getter,
		storer: storer,
	}
}

// ServeHTTP processes category creation
func (h *postCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// This category api model will hold all
	// the params from the request
	categoryAPI := category{}
	err := categoryAPI.fromRequest(r)
	if err != nil {
		if err == ErrInvalidBody {
			log.Info(err)
			api.RenderInvalidInput(w, r, "", ErrInvalidBody.Error())
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	// This category model contains all the values
	// for the category we want to create
	modelCategory := categoryAPI.toModel()

	if err := h.storer.Store(&modelCategory); err != nil {
		api.RenderInternalServerError(w, r, fmt.Errorf("postCategory store error: %s", err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"category_id": modelCategory.CategoryID})
}

package category

import (
	"github.com/go-chi/render"
	"github.com/gregbiv/news-api/pkg/api"
	"github.com/gregbiv/news-api/pkg/context"
	"github.com/gregbiv/news-api/pkg/model"
	storageCategory "github.com/gregbiv/news-api/pkg/storage/category"
	"net/http"
)

type (
	categoryUpdateHandler struct {
		getter       storageCategory.Getter
		updater      storageCategory.Updater
		urlExtractor api.URLExtractor
	}
)

// NewCategoryUpdateHandler init and returns an instance of categoryUpdateHandler
func NewCategoryUpdateHandler(
	getter storageCategory.Getter,
	updater storageCategory.Updater,
	urlExtractor api.URLExtractor,
) http.Handler {
	return &categoryUpdateHandler{
		getter,
		updater,
		urlExtractor,
	}
}

func (h *categoryUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	categoryID, err := h.urlExtractor.UUIDFromRoute(r, "category_id")
	if err != nil {
		api.NotFound(w, r)
		return
	}

	categoryAPI := category{}
	categoryAPI.CategoryID = categoryID

	err = categoryAPI.fromRequest(r)
	if err == ErrInvalidBody {
		context.Logger(r.Context()).Info(err)
		api.RenderInvalidInput(w, r, "", ErrInvalidBody.Error())
		return
	}
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	dbCategory, err := h.getter.GetCategoryByID(categoryAPI.CategoryID.String())
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	// This categoryAPI model contains all the new values
	// for the categoryAPI we want to update
	modelCategory := categoryAPI.toModel()

	h.setCategory(&modelCategory, dbCategory)

	err = h.updater.Update(dbCategory)
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	response := category{}
	err = response.fromDB(dbCategory)
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func (h *categoryUpdateHandler) setCategory(source, target *model.Category) {
	target.CategoryID = source.CategoryID
	target.Name = source.Name
	target.Title = source.Title
}

package docs

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gregbiv/news-api/pkg/api"
	"github.com/gregbiv/news-api/pkg/assets/docs"
)

// Docs registers sub routes for all documentation assets
func Docs(r chi.Router) {
	for _, file := range docs.AssetNames() {
		asset, err := docs.AssetInfo(file)
		if err != nil {
			panic(err)
		}

		r.Get("/"+file, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			content, err := docs.Asset(asset.Name())
			if err != nil {
				api.NotFound(w, r)
				return
			}

			http.ServeContent(w, r, asset.Name(), asset.ModTime(), bytes.NewReader(content))
		}))
	}
}

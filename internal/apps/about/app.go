package about

import (
	"chia-goths/internal/apps"
	"embed"

	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed assets/*
var embeddedAssetsFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

//go:embed articles/*
var articlesFS embed.FS

type App struct{}

func (a App) Init(config *apps.AppConfig) {
	c := config.Router
	renderer := config.Renderer

	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.RenderHTML(r, w, "index", nil)
	})
	c.Get("/articles/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		content, err := articlesFS.ReadFile("articles/" + name + ".md")
		if err != nil {
			renderer.RenderHTML(r, w, "404", nil)
			return
		}

		renderer.RenderHTML(r, w, "article", map[string]interface{}{
			"Content": content,
		})
	})
}

func (a App) GetAssetsFS() apps.AssetsFS {
	return apps.AssetsFS{
		EmbeddedFS:   embeddedAssetsFS,
		RelativePath: "assets",
	}
}

func (a App) GetTemplatesEmbedFS() apps.TemplatesFS {
	return apps.TemplatesFS{
		EmbeddedFS:   templatesFS,
		RelativePath: "templates",
	}
}

func (a App) GetAppPath() string {
	return "internal/apps/about"
}

package plugins

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/wellington-vell/gorpc"
)

type ScalarPlugin struct {
	app         *gorpc.GORPC
	openAPIPath string
	uiPath      string
	title       string
	theme       string
	layout      string
	config      map[string]any
}

type ScalarPluginOptions struct {
	UIPath string
	Title  string
	Theme  string
	Layout string
	Config map[string]any
}

func NewScalarPlugin(options ...*ScalarPluginOptions) *ScalarPlugin {
	uiPath := "/scalar"
	title := "API Documentation - Scalar"
	theme := "purple"
	layout := "modern"
	config := make(map[string]any)

	if len(options) > 0 && options[0] != nil {
		opts := options[0]
		if opts.UIPath != "" {
			uiPath = opts.UIPath
		}
		if opts.Title != "" {
			title = opts.Title
		}
		if opts.Theme != "" {
			theme = opts.Theme
		}
		if opts.Layout != "" {
			layout = opts.Layout
		}
		if opts.Config != nil {
			config = opts.Config
		}
	}

	return &ScalarPlugin{
		openAPIPath: "/openapi.json",
		uiPath:      uiPath,
		title:       title,
		theme:       theme,
		layout:      layout,
		config:      config,
	}
}

func (p *ScalarPlugin) Name() string {
	return "scalar"
}

func (p *ScalarPlugin) Register(app *gorpc.GORPC) error {
	p.app = app
	if !app.HasPlugin("openapi") {
		app.Plugin(NewOpenAPIPlugin())
	}
	return nil
}

func (p *ScalarPlugin) Routes() map[string]http.Handler {
	return map[string]http.Handler{
		p.uiPath: http.HandlerFunc(p.serveUI),
	}
}

func (p *ScalarPlugin) serveUI(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]any)
	maps.Copy(config, p.config)
	config["theme"] = p.theme
	config["layout"] = p.layout

	configJSON, err := json.Marshal(config)
	if err != nil {
		configJSON = []byte(`{"theme": "purple", "layout": "modern"}`)
	}

	html := fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <title>%s</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      type="application/json"
      data-configuration='%s'
      data-url="%s"
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest"></script>
  </body>
</html>`, p.title, string(configJSON), p.openAPIPath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

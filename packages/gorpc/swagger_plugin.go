package gorpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SwaggerPlugin struct {
	app         *GORPC
	openAPIPath string
	uiPath      string
	title       string
	layout      string
	deepLinking bool
	config      map[string]interface{}
}

type SwaggerPluginOptions struct {
	UIPath     string
	Title      string
	Layout     string
	DeepLinking *bool
	Config     map[string]interface{}
}

func NewSwaggerPlugin(options ...*SwaggerPluginOptions) *SwaggerPlugin {
	uiPath := "/swagger/"
	title := "Swagger UI"
	layout := "StandaloneLayout"
	deepLinking := true
	config := make(map[string]interface{})

	if len(options) > 0 && options[0] != nil {
		opts := options[0]
		if opts.UIPath != "" {
			uiPath = opts.UIPath
			if len(uiPath) > 0 && uiPath[len(uiPath)-1] != '/' {
				uiPath = uiPath + "/"
			}
		}
		if opts.Title != "" {
			title = opts.Title
		}
		if opts.Layout != "" {
			layout = opts.Layout
		}
		if opts.DeepLinking != nil {
			deepLinking = *opts.DeepLinking
		}
		if opts.Config != nil {
			config = opts.Config
		}
	}

	return &SwaggerPlugin{
		openAPIPath: "/openapi.json",
		uiPath:      uiPath,
		title:       title,
		layout:       layout,
		deepLinking:  deepLinking,
		config:       config,
	}
}

func (p *SwaggerPlugin) Name() string {
	return "swagger"
}

func (p *SwaggerPlugin) Register(app *GORPC) error {
	p.app = app
	if !app.hasPlugin("openapi") {
		app.Plugin(NewOpenAPIPlugin())
	}
	return nil
}

func (p *SwaggerPlugin) Routes() map[string]http.Handler {
	basePath := p.uiPath
	if len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
		basePath = basePath[:len(basePath)-1]
	}

	return map[string]http.Handler{
		basePath:                http.HandlerFunc(p.serveUI),
		p.uiPath:                http.HandlerFunc(p.serveUI),
		p.uiPath + "index.html": http.HandlerFunc(p.serveUI),
	}
}

func (p *SwaggerPlugin) serveUI(w http.ResponseWriter, r *http.Request) {
	basePath := p.uiPath
	if len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
		basePath = basePath[:len(basePath)-1]
	}
	if r.URL.Path == basePath {
		http.Redirect(w, r, p.uiPath, http.StatusMovedPermanently)
		return
	}

	swaggerConfig := map[string]interface{}{
		"url":         p.openAPIPath,
		"dom_id":      "#swagger-ui",
		"layout":      p.layout,
		"deepLinking": p.deepLinking,
	}
	for k, v := range p.config {
		if k != "presets" {
			swaggerConfig[k] = v
		}
	}
	configJSON, err := json.Marshal(swaggerConfig)
	if err != nil {
		configJSON = []byte(fmt.Sprintf(`{"url": "%s", "dom_id": "#swagger-ui", "layout": "%s", "deepLinking": %t}`, p.openAPIPath, p.layout, p.deepLinking))
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>%s</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js" crossorigin></script>
    <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-standalone-preset.js" crossorigin></script>
    <script>
      window.onload = function() {
        var config = %s;
        config.presets = [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ];
        window.ui = SwaggerUIBundle(config);
      };
    </script>
  </body>
</html>`, p.title, string(configJSON))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

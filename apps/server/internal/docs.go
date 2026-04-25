package internal

import (
	"net/http"
	"os"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

func Docs(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		http.Error(w, "Swagger doc not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		return
	}
}

func SwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: "/spec.json",
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		return
	}
}

func ScalarUI(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		http.Error(w, "Swagger doc not found", http.StatusNotFound)
		return
	}

	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: string(data),
	})
	if err != nil {
		http.Error(w, "Failed to render Scalar UI", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		return
	}
}

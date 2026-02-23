package plugins

import (
	"net/http"

	"github.com/wellington-vell/gorpc"
	"github.com/wellington-vell/gorpc/openapi"
)

type OpenAPIPlugin struct {
	app  *gorpc.GORPC
	path string
}

func NewOpenAPIPlugin() *OpenAPIPlugin {
	return &OpenAPIPlugin{
		path: "/openapi.json",
	}
}

func (p *OpenAPIPlugin) Name() string {
	return "openapi"
}

func (p *OpenAPIPlugin) Register(app *gorpc.GORPC) error {
	p.app = app
	return nil
}

func (p *OpenAPIPlugin) Routes() map[string]http.Handler {
	return map[string]http.Handler{
		p.path: http.HandlerFunc(p.serveOpenAPI),
	}
}

func (p *OpenAPIPlugin) serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	procedures := p.app.GetAllProcedures()

	openAPIProcedures := make([]openapi.ProcedureInfo, len(procedures))
	for i, proc := range procedures {
		openAPIProcedures[i] = openapi.ProcedureInfo{
			Path:       proc.Path,
			Method:     proc.Method,
			Route:      proc.Route,
			Meta:       proc.Meta,
			Tags:       proc.Tags,
			InputType:  proc.InputType,
			OutputType: proc.OutputType,
			ErrorCodes: proc.ErrorCodes,
			PathParams: proc.PathParams,
		}
	}

	handler := openapi.Handler(openAPIProcedures)
	handler.ServeHTTP(w, r)
}

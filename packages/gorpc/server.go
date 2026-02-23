package gorpc

import (
	"log"
	"net/http"
)

func (g *GORPC) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	for _, plugin := range g.pluginRegistry.Get() {
		routes := plugin.Routes()
		for path, handler := range routes {
			mux.Handle(path, g.wrapCORS(handler))
		}
	}
	mux.Handle("/", g.wrapCORS(g.router))

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, mux)
}

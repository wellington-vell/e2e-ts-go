package gorpc

import (
	"net/http"
)

// Plugin is the interface that all gorpc plugins must implement.
// Plugins can register routes, middleware, or modify the application behavior.
type Plugin interface {
	// Name returns a unique identifier for the plugin
	Name() string

	// Register is called when the plugin is registered with the GORPC instance.
	// It receives the GORPC instance and can use it to register routes, access routers, etc.
	Register(app *GORPC) error

	// Routes returns HTTP handlers that the plugin wants to mount.
	// The map key is the path prefix, and the value is the http.Handler to mount at that path.
	// If a plugin doesn't need to mount routes, it can return nil or an empty map.
	Routes() map[string]http.Handler
}

// pluginRegistry maintains a list of registered plugins and provides methods
// for adding and retrieving them. It is used internally by GORPC to manage
// plugin lifecycle and route registration.
type pluginRegistry struct {
	plugins []Plugin
}

func newPluginRegistry() *pluginRegistry {
	return &pluginRegistry{
		plugins: make([]Plugin, 0),
	}
}

func (pr *pluginRegistry) Register(plugin Plugin) {
	if plugin == nil {
		return
	}
	pr.plugins = append(pr.plugins, plugin)
}

func (pr *pluginRegistry) Get() []Plugin {
	return pr.plugins
}

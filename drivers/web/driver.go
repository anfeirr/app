// Package web is the driver to be used for web applications.
// It is build on the top of GopherJS.
package web

import (
	"net/http"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is an app.Driver implementation for web.
type Driver struct {
	// The URL of the component to load when a navigating on the website root.
	URL string

	// The URL of the component to load when a 404 errors occurs.
	// Default is /web.NotFound
	NotFoundURL string

	// The app icon name.
	Icon string

	// The server used to save request.
	// Default is a server that listens on port 7042.
	Server *http.Server

	// OnServerRun is called when the web server is running.
	// http.Handler overrides should be performed here.
	OnServerRun func()

	ui          chan func()
	factory     *app.Factory
	events      *app.EventRegistry
	elems       *core.ElemDB
	page        app.Page
	stop        func()
	fileHandler http.Handler
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "web"
}

// Name satisfies the app.Driver interface.
func (d *Driver) Name() string {
	return "Web"
}

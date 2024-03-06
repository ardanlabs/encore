// Package all binds all the routes into the specified app.
package all

import (
	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/usergrp"
	"github.com/ardanlabs/encore/business/web/v1/mux"
	"github.com/ardanlabs/encore/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	usergrp.Routes(app, usergrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})
}

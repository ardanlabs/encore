// Package delegate provides the ability to make function calls between
// different core packages when an import is not possible.
package delegate

import (
	"context"

	"encore.dev/rlog"
)

// These types are just for documentation so we know what keys go
// where in the map.
type (
	domain string
	action string
)

// Delegate manages the set of functions to be called by core
// packages when an import is not possible.
type Delegate struct {
	funcs map[domain]map[action][]Func
}

// New constructs a delegate for indirect api access.
func New() *Delegate {
	return &Delegate{
		funcs: make(map[domain]map[action][]Func),
	}
}

// Register adds a function to be called for a specified domain and action.
func (d *Delegate) Register(domainType string, actionType string, fn Func) {
	aMap, ok := d.funcs[domain(domainType)]
	if !ok {
		aMap = make(map[action][]Func)
		d.funcs[domain(domainType)] = aMap
	}

	funcs := aMap[action(actionType)]
	funcs = append(funcs, fn)
	aMap[action(actionType)] = funcs
}

// Call executes all functions registered for the specified domain and
// action. These functions are executed synchronously on the G making the call.
func (d *Delegate) Call(ctx context.Context, data Data) error {
	rlog.Info("delegate call", "status", "started", "domain", data.Domain, "action", data.Action, "params", data.RawParams)
	defer rlog.Info("delegate call", "status", "completed")

	if dMap, ok := d.funcs[domain(data.Domain)]; ok {
		if funcs, ok := dMap[action(data.Action)]; ok {
			for _, fn := range funcs {
				rlog.Info("delegate call", "status", "sending")

				if err := fn(ctx, data); err != nil {
					rlog.Error("delegate call", "msg", err)
				}
			}
		}
	}

	return nil
}

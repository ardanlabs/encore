package auth

import (
	"context"
	"runtime"

	"encore.dev/rlog"
)

// Service represents the encore service application.
//
//encore:service
type Service struct {
	log rlog.Ctx
}

// NewService is called to create a new encore Service.
func NewService(log rlog.Ctx) (*Service, error) {
	s := Service{
		log: log,
	}

	return &s, nil
}

// Shutdown implements a function that will be called by encore when the service
// is signaled to shutdown.
func (s *Service) Shutdown(force context.Context) {
	defer s.log.Info("shutdown", "status", "shutdown complete")
}

// =============================================================================

// initService is called by Encore to initialize the service.
//
//lint:ignore U1000 "called by encore"
func initService() (*Service, error) {
	log := rlog.With("service", "auth")

	err := startup(log)
	if err != nil {
		return nil, err
	}

	return NewService(log)
}

func startup(log rlog.Ctx) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info("initService", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	return nil
}

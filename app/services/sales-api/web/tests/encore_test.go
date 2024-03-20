package tests

import (
	"context"

	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/usergrp"
	"github.com/ardanlabs/encore/app/services/sales-api/web/service"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/database"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/jmoiron/sqlx"
)

//encore:service
type Service struct {
	Log     *logger.Logger
	DB      *sqlx.DB
	Auth    *auth.Auth
	UsrCore *user.Core
	UsrGrp  *usergrp.Handlers
}

// initService is called by Encore to initialize the service.
func initService() (*Service, error) {
	s, err := service.New(database.EDB)
	if err != nil {
		return nil, err
	}

	es := Service{
		Log:     s.Log,
		DB:      s.DB,
		Auth:    s.Auth,
		UsrCore: s.UsrCore,
		UsrGrp:  s.UsrGrp,
	}

	return &es, nil
}

// Shutdown implements a function that will be called by encore when the service
// is signaled to shutdown.
func (s *Service) Shutdown(force context.Context) {
	defer s.Log.Info(force, "shutdown", "status", "shutdown complete")

	s.Log.Info(force, "shutdown", "status", "stopping database support")
	s.DB.Close()
}

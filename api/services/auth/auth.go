// Package auth represent the encore application.
package auth

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"encore.dev"
	esqldb "encore.dev/storage/sqldb"
	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/encore/app/sdk/auth"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/sdk/delegate"
	"github.com/ardanlabs/encore/business/sdk/sqldb"
	"github.com/ardanlabs/encore/foundation/keystore"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Represents the database this service will use. The name has to be a literal
// string.
var appDB = esqldb.Named("app")

// Represents the secrets for this service.
var secrets struct {
	KeyID  string
	KeyPEM string
}

// =============================================================================

// Service represents the encore service application.
//
//encore:service
type Service struct {
	log     *logger.Logger
	db      *sqlx.DB
	auth    *auth.Auth
	userBus *userbus.Business
}

// NewService is called to create a new encore Service.
func NewService(log *logger.Logger, db *sqlx.DB, ath *auth.Auth) (*Service, error) {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, userdb.NewStore(log, db))

	s := Service{
		log:     log,
		db:      db,
		auth:    ath,
		userBus: userBus,
	}

	return &s, nil
}

// Shutdown implements a function that will be called by encore when the service
// is signaled to shutdown.
func (s *Service) Shutdown(force context.Context) {
	ctx := context.Background()

	defer s.log.Info(ctx, "shutdown", "status", "shutdown complete")

	s.log.Info(ctx, "shutdown", "status", "stopping database support")
	s.db.Close()
}

// =============================================================================

// initService is called by Encore to initialize the service.
//
//lint:ignore U1000 "called by encore"
func initService() (*Service, error) {
	log := logger.New("auth")

	db, auth, err := startup(log)
	if err != nil {
		return nil, err
	}

	return NewService(log, db, auth)
}

func startup(log *logger.Logger) (*sqlx.DB, *auth.Auth, error) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "initService", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Auth struct {
			ActiveKID string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
			Issuer    string `conf:"default:service project"`
		}
		DB struct {
			MaxIdleConns int `conf:"default:0"`
			MaxOpenConns int `conf:"default:0"`
		}
	}{
		Version: conf.Version{
			Build: encore.Meta().Environment.Name,
			Desc:  "Auth",
		},
	}

	const prefix = "AUTH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "initService", "environment", encore.Meta().Environment.Name)

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "initService", "config", out)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "initService", "status", "initializing database support")

	db, err := sqldb.Open(sqldb.Config{
		EDB:          appDB,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to db: %w", err)
	}

	// -------------------------------------------------------------------------
	// Auth Support

	log.Info(ctx, "initService", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.

	ks := keystore.New()
	if err := ks.LoadKey(secrets.KeyID, secrets.KeyPEM); err != nil {
		return nil, nil, fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
		Issuer:    cfg.Auth.Issuer,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("constructing auth: %w", err)
	}

	return db, auth, nil
}

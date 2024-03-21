// Package encore represent the encore application.
package encore

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"encore.dev/rlog"
	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/homegrp"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/productgrp"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/trangrp"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/usergrp"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/vproductgrp"
	"github.com/ardanlabs/encore/business/core/crud/delegate"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/home/stores/homedb"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/product/stores/productdb"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/core/crud/user/stores/userdb"
	"github.com/ardanlabs/encore/business/core/views/vproduct"
	"github.com/ardanlabs/encore/business/core/views/vproduct/stores/vproductdb"
	"github.com/ardanlabs/encore/business/data/appdb"
	"github.com/ardanlabs/encore/business/data/appdb/migrate"
	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/debug"
	"github.com/ardanlabs/encore/business/web/metrics"
	"github.com/ardanlabs/encore/foundation/keystore"
	"github.com/jmoiron/sqlx"
)

var build = "develop"

//encore:service
type Service struct {
	mtrcs   *metrics.Values
	db      *sqlx.DB
	auth    *auth.Auth
	usrCore *user.Core
	prdCore *product.Core
	hmeCore *home.Core
	usrGrp  *usergrp.Handlers
	prdGrp  *productgrp.Handlers
	hmeGrp  *homegrp.Handlers
	trnGrp  *trangrp.Handlers
	vprdGrp *vproductgrp.Handlers
	debug   http.Handler
}

// initService is called by Encore to initialize the service.
//
//lint:ignore U1000 "called by encore"
func initService() (*Service, error) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	rlog.Info("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
			Issuer     string `conf:"default:service project"`
		}
		DB struct {
			MaxIdleConns int `conf:"default:2"`
			MaxOpenConns int `conf:"default:0"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Service Project",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, err
		}
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	rlog.Info("starting service", "version", build)
	defer rlog.Info("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, fmt.Errorf("generating config for output: %w", err)
	}
	rlog.Info("startup", "config", out)

	expvar.NewString("build").Set(build)

	// -------------------------------------------------------------------------
	// Database Support

	rlog.Info("startup", "status", "initializing database support")

	db, err := sqldb.Open(sqldb.Config{
		EDB:          appdb.AppDB,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	// TODO: I don't like this here because it's more of an ops thing, but
	// for now I will leave it as I learn more.
	migrate.Seed(ctx, db)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	rlog.Info("startup", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return nil, fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		DB:        db,
		KeyLookup: ks,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return nil, fmt.Errorf("constructing auth: %w", err)
	}

	usrCore := user.NewCore(delegate.New(), userdb.NewStore(db))
	prdCore := product.NewCore(usrCore, delegate.New(), productdb.NewStore(db))
	hmeCore := home.NewCore(usrCore, delegate.New(), homedb.NewStore(db))
	vprdCore := vproduct.NewCore(vproductdb.NewStore(db))

	s := Service{
		mtrcs:   newMetrics(),
		db:      db,
		auth:    auth,
		usrCore: usrCore,
		prdCore: prdCore,
		hmeCore: hmeCore,
		usrGrp:  usergrp.New(usrCore, auth),
		prdGrp:  productgrp.New(prdCore),
		hmeGrp:  homegrp.New(hmeCore),
		trnGrp:  trangrp.New(usrCore, prdCore),
		vprdGrp: vproductgrp.New(vprdCore),

		debug: debug.Mux(),
	}

	return &s, nil
}

// Shutdown implements a function that will be called by encore when the service
// is signaled to shutdown.
func (s *Service) Shutdown(force context.Context) {
	defer rlog.Info("shutdown", "status", "shutdown complete")

	rlog.Info("shutdown", "status", "stopping database support")
	s.db.Close()
}

// Fallback is called for the debug enpoints.
//
//encore:api public raw path=/!fallback
func (s *Service) Fallback(w http.ResponseWriter, req *http.Request) {
	s.debug.ServeHTTP(w, req)
}

package encore

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	edb "encore.dev/storage/sqldb"
	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/build/all"
	"github.com/ardanlabs/encore/business/core/crud/delegate"
	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/business/web/v1/debug"
	"github.com/ardanlabs/encore/business/web/v1/mux"
	"github.com/ardanlabs/encore/foundation/keystore"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/ardanlabs/encore/foundation/web"
	"github.com/jmoiron/sqlx"
)

var build = "develop"

var ebdDB = edb.NewDatabase("url", edb.DatabaseConfig{
	Migrations: "./v1/migrations",
})

//encore:service
type Service struct {
	log             *logger.Logger
	db              *sqlx.DB
	api             *http.Server
	shutdownTimeout time.Duration
}

// initService is called by Encore to initialize the service.
func initService() (*Service, error) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Logging Starting

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES-API", traceIDFn, events)

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:4000"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
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

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(build)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support")

	db, err := sqldb.Open(sqldb.Config{
		EDB:          ebdDB,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	seedDatabase(ctx, db)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return nil, fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return nil, fmt.Errorf("constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:    build,
		Shutdown: shutdown,
		Log:      log,
		Delegate: delegate.New(log),
		Auth:     auth,
		DB:       db,
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfgMux, all.Routes(), mux.WithCORS(cfg.Web.CORSAllowedOrigins)),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	s := Service{
		log:             log,
		db:              db,
		api:             &api,
		shutdownTimeout: cfg.Web.ShutdownTimeout,
	}

	return &s, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.log.Info(force, "shutdown", "status", "shutdown complete")

	ctx, cancel := context.WithTimeout(force, s.shutdownTimeout)
	defer cancel()

	if err := s.api.Shutdown(ctx); err != nil {
		s.api.Close()
		s.log.Info(force, "shutdown", "could not stop server gracefully: %w", err)
	}

	s.log.Info(force, "shutdown", "status", "stopping database support")
	s.db.Close()
}

// Route all requests to the existing HTTP router if no other endpoint matches.
//
//encore:api public raw path=/!fallback
func (s *Service) Fallback(w http.ResponseWriter, req *http.Request) {
	s.api.Handler.ServeHTTP(w, req)
}

//go:embed v1/seeds/seed.sql
var seedDoc string

func seedDatabase(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(seedDoc); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

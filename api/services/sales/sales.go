// Package sales represent the encore application.
package sales

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"encore.dev"
	esqldb "encore.dev/storage/sqldb"
	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/encore/app/domain/homeapp"
	"github.com/ardanlabs/encore/app/domain/productapp"
	"github.com/ardanlabs/encore/app/domain/tranapp"
	"github.com/ardanlabs/encore/app/domain/userapp"
	"github.com/ardanlabs/encore/app/domain/vproductapp"
	"github.com/ardanlabs/encore/app/sdk/debug"
	"github.com/ardanlabs/encore/app/sdk/metrics"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/homebus/stores/homedb"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/productbus/stores/productdb"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/ardanlabs/encore/business/domain/vproductbus/stores/vproductdb"
	"github.com/ardanlabs/encore/business/sdk/appdb/migrate"
	"github.com/ardanlabs/encore/business/sdk/delegate"
	"github.com/ardanlabs/encore/business/sdk/sqldb"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Represents the database this service will use. The name has to be a literal
// string.
var appDB = esqldb.Named("app")

// =============================================================================

// Service represents the encore service application.
//
//encore:service
type Service struct {
	log   *logger.Logger
	mtrcs *metrics.Values
	db    *sqlx.DB
	debug http.Handler
	appDomain
	busDomain
}

// NewService is called to create a new encore Service.
func NewService(log *logger.Logger, db *sqlx.DB) (*Service, error) {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, userdb.NewStore(log, db))
	productBus := productbus.NewBusiness(log, userBus, delegate, productdb.NewStore(log, db))
	homeBus := homebus.NewBusiness(log, userBus, delegate, homedb.NewStore(log, db))
	vproductBus := vproductbus.NewBusiness(vproductdb.NewStore(log, db))

	s := Service{
		log:   log,
		mtrcs: newMetrics(),
		db:    db,
		debug: debug.Mux(),
		appDomain: appDomain{
			userApp:     userapp.NewApp(userBus),
			productApp:  productapp.NewApp(productBus),
			homeApp:     homeapp.NewApp(homeBus),
			tranApp:     tranapp.NewApp(userBus, productBus),
			vproductApp: vproductapp.NewApp(vproductBus),
		},
		busDomain: busDomain{
			delegate:   delegate,
			userBus:    userBus,
			productBus: productBus,
			homeBus:    homeBus,
		},
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
	log := logger.New("sales")

	db, err := startup(log)
	if err != nil {
		return nil, err
	}

	return NewService(log, db)
}

func startup(log *logger.Logger) (*sqlx.DB, error) {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "initService", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		DB struct {
			MaxIdleConns int `conf:"default:0"`
			MaxOpenConns int `conf:"default:0"`
		}
	}{
		Version: conf.Version{
			Build: encore.Meta().Environment.Name,
			Desc:  "Sales",
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

	log.Info(ctx, "initService", "environment", encore.Meta().Environment.Name)

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, fmt.Errorf("generating config for output: %w", err)
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
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	if err := migrate.Seed(context.Background(), db); err != nil {
		return nil, fmt.Errorf("seeding the db: %w", err)
	}

	return db, nil
}

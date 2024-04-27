// Package homedb contains home related CRUD functionality.
package homedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"encore.dev/rlog"
	"github.com/ardanlabs/encore/business/api/order"
	"github.com/ardanlabs/encore/business/api/sqldb"
	"github.com/ardanlabs/encore/business/api/transaction"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for home database access.
type Store struct {
	log rlog.Ctx
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log rlog.Ctx, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (homebus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		db: ec,
	}

	return &store, nil
}

// Create inserts a new home into the database.
func (s *Store) Create(ctx context.Context, hme homebus.Home) error {
	const q = `
    INSERT INTO homes
        (home_id, user_id, type, address_1, address_2, zip_code, city, state, country, date_created, date_updated)
    VALUES
        (:home_id, :user_id, :type, :address_1, :address_2, :zip_code, :city, :state, :country, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBHome(hme)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a home from the database.
func (s *Store) Delete(ctx context.Context, hme homebus.Home) error {
	data := struct {
		ID string `db:"home_id"`
	}{
		ID: hme.ID.String(),
	}

	const q = `
    DELETE FROM
	    homes
	WHERE
	  	home_id = :home_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a home document in the database.
func (s *Store) Update(ctx context.Context, hme homebus.Home) error {
	const q = `
    UPDATE
        homes
    SET
        "address_1"     = :address_1,
        "address_2"     = :address_2,
        "zip_code"      = :zip_code,
        "city"          = :city,
        "state"         = :state,
        "country"       = :country,
        "type"          = :type,
        "date_updated"  = :date_updated
    WHERE
        home_id = :home_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBHome(hme)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing homes from the database.
func (s *Store) Query(ctx context.Context, filter homebus.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]homebus.Home, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
    SELECT
	    home_id, user_id, type, address_1, address_2, zip_code, city, state, country, date_created, date_updated
	FROM
	  	homes`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbHmes []dbHome
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbHmes); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	hmes, err := toCoreHomeSlice(dbHmes)
	if err != nil {
		return nil, err
	}

	return hmes, nil
}

// Count returns the total number of homes in the DB.
func (s *Store) Count(ctx context.Context, filter homebus.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        homes`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified home from the database.
func (s *Store) QueryByID(ctx context.Context, homeID uuid.UUID) (homebus.Home, error) {
	data := struct {
		ID string `db:"home_id"`
	}{
		ID: homeID.String(),
	}

	const q = `
    SELECT
	  	home_id, user_id, type, address_1, address_2, zip_code, city, state, country, date_created, date_updated
    FROM
        homes
    WHERE
        home_id = :home_id`

	var dbHme dbHome
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbHme); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return homebus.Home{}, fmt.Errorf("db: %w", homebus.ErrNotFound)
		}
		return homebus.Home{}, fmt.Errorf("db: %w", err)
	}

	return toCoreHome(dbHme)
}

// QueryByUserID gets the specified home from the database by user id.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]homebus.Home, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    home_id, user_id, type, address_1, address_2, zip_code, city, state, country, date_created, date_updated
	FROM
		homes
	WHERE
		user_id = :user_id`

	var dbHmes []dbHome
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbHmes); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toCoreHomeSlice(dbHmes)
}

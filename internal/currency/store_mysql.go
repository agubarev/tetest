package currency

import (
	"context"
	"database/sql"

	"github.com/agubarev/tetest/util/guard"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type defaultMySQLStore struct {
	connection *dbr.Connection
}

func NewDefaultMySQLStore(connection *dbr.Connection) (Store, error) {
	if connection == nil {
		return nil, ErrNilDatabase
	}

	s := &defaultMySQLStore{
		connection: connection,
	}

	return s, nil
}

func (s *defaultMySQLStore) oneByQuery(ctx context.Context, q string, args ...interface{}) (c Currency, err error) {
	err = s.connection.NewSession(nil).
		SelectBySql(q, args...).
		LoadOneContext(ctx, &c)

	if err != nil {
		if err == sql.ErrNoRows {
			return c, ErrCurrencyNotFound
		}

		return c, err
	}

	return c, nil
}

func (s *defaultMySQLStore) manyByQuery(ctx context.Context, q string, args ...interface{}) (items []Currency, err error) {
	items = make([]Currency, 0)

	_, err = s.connection.NewSession(nil).
		SelectBySql(q, args...).
		LoadContext(ctx, &items)

	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}

		return nil, err
	}

	return items, nil
}

// CreateCurrency creates a new entry in the storage backend
func (s *defaultMySQLStore) BulkCreateCurrency(ctx context.Context, items []Currency) (_ []Currency, err error) {
	currencyLen := len(items)

	// there must be something first
	if currencyLen == 0 {
		return nil, ErrNoData
	}

	tx, err := s.connection.NewSession(nil).Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize database transaction")
	}
	defer tx.RollbackUnlessCommitted()

	//---------------------------------------------------------------------------
	// building the bulk statement
	//---------------------------------------------------------------------------
	stmt := tx.InsertInto("currency").Columns(guard.DBColumnsFrom(items[0])...)

	// validating each c individually
	for i := range items {
		if err := items[i].validate(); err != nil {
			return nil, err
		}

		// adding value to the batch
		stmt = stmt.Record(items[i])
	}

	// executing the batch
	result, err := stmt.ExecContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "bulk insert failed")
	}

	// returned ID belongs to the first c created
	firstNewID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit database transaction")
	}

	// distributing new IDs in their sequential order
	for i := range items {
		items[i].ID = int(firstNewID)
		firstNewID++
	}

	return items, nil
}

func (s *defaultMySQLStore) FetchCurrencyByID(ctx context.Context, id int) (c Currency, err error) {
	return s.oneByQuery(ctx, "SELECT * FROM `c` WHERE id = ? LIMIT 1", id)
}

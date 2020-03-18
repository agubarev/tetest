package currency

import (
	"context"
	"database/sql"
	"time"

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

func (s *defaultMySQLStore) BulkCreate(ctx context.Context, cs []Currency) (_ []Currency, err error) {
	currencyLen := len(cs)

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
	stmt := tx.InsertInto("currency").Columns(guard.DBColumnsFrom(cs[0])...)

	// validating each c individually
	for i := range cs {
		if err := cs[i].validate(); err != nil {
			return nil, err
		}

		// adding value to the batch
		stmt = stmt.Record(&cs[i])
	}

	// executing the batch
	if _, err = stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrap(err, "bulk insert failed")
	}

	// NOTE: since this is a very simple test case, thus
	// ignoring execution result and not assigning
	// keys (because in this case ID is the currency name)

	return cs, nil
}

func (s *defaultMySQLStore) ByID(ctx context.Context, id string) (cs []Currency, err error) {
	return s.manyByQuery(ctx, "SELECT * FROM `currency` WHERE id = ?", id)
}

func (s *defaultMySQLStore) ByDate(ctx context.Context, date time.Time) (cs []Currency, err error) {
	return s.manyByQuery(ctx, "SELECT * FROM `currency` WHERE `valid_date` = ? LIMIT 1", dbr.Expr("DATE(?)", time.Time{}))
}

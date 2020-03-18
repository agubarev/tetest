package currency

import (
	"context"
	"database/sql"
	"time"

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
	// NOTE: because I want currency values to be updated on repetitive
	// insert attempts (i.e. multiple import runs or external changes),
	// so I'm using a prepared statement, otherwise I'd simply go for the following:
	// stmt := tx.InsertInto("currency").Columns(guard.DBColumnsFrom(&cs[0])...)

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO currency(id, value, pub_date, created_at)
		VALUES(?, ?, ?, NOW()) ON DUPLICATE KEY UPDATE value = ?, updated_at = NOW()`)

	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statement")
	}

	// statement must be closed afterwards
	defer func() {
		if err = stmt.Close(); err != nil {
			err = errors.Wrap(err, "failed to close prepared statement")
		}
	}()

	// validating each c individually
	for i := range cs {
		if err := cs[i].validate(); err != nil {
			return nil, err
		}

		c := cs[i]

		// NOTE: ignoring execution result
		if _, err = stmt.ExecContext(ctx, c.ID, c.Value, c.PubDate, c.Value); err != nil {
			return nil, errors.Wrap(err, "failed to execute statement")
		}
	}

	// committing
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit database transaction")
	}

	// NOTE: since this is a very simple test case, thus
	// ignoring execution result and not assigning
	// keys (because in this case ID is the currency name)

	return cs, nil
}

func (s *defaultMySQLStore) ByID(ctx context.Context, id string) (cs []Currency, err error) {
	return s.manyByQuery(ctx, "SELECT * FROM `currency` WHERE id = ? ORDER BY created_at DESC, updated_at DESC", id)
}

func (s *defaultMySQLStore) ByDate(ctx context.Context, date time.Time) (cs []Currency, err error) {
	return s.manyByQuery(ctx, "SELECT * FROM `currency` WHERE `pub_date` = ?", dbr.Expr("DATE(?)", date.Local()))
}

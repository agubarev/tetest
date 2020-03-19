package currency

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMySQLStore_BulkCreate(t *testing.T) {
	a := assert.New(t)

	db, mock, err := sqlmock.New()
	a.NoError(err)
	a.NotNil(db)
	a.NotNil(mock)
	defer db.Close()

	store, err := NewDefaultMySQLStore(&dbr.Connection{
		DB:            db,
		Dialect:       dialect.MySQL,
		EventReceiver: nil,
	})

	testdata := []Currency{
		{ID: "LVL", Value: 1.00, PubDate: dbr.NewNullTime(time.Now())},
		{ID: "EUR", Value: 2.00, PubDate: dbr.NewNullTime(time.Now())},
		{ID: "USD", Value: 3.00, PubDate: dbr.NewNullTime(time.Now())},
	}

	mock.ExpectBegin()

	// assigning to variable for re-use due to multiple exec calls,
	stmt := mock.ExpectPrepare("INSERT INTO currency")

	stmt.ExpectExec().
		WithArgs(testdata[0].ID, testdata[0].Value, testdata[0].PubDate, testdata[0].Value).
		WillReturnResult(sqlmock.NewResult(0, 1))

	stmt.ExpectExec().
		WithArgs(testdata[1].ID, testdata[1].Value, testdata[1].PubDate, testdata[1].Value).
		WillReturnResult(sqlmock.NewResult(0, 1))

	stmt.ExpectExec().
		WithArgs(testdata[2].ID, testdata[2].Value, testdata[2].PubDate, testdata[2].Value).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	cs, err := store.BulkCreate(context.Background(), testdata)
	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, len(testdata))

	a.NoError(mock.ExpectationsWereMet())
}

func TestDefaultMySQLStore_AllLatest(t *testing.T) {
	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	a.NoError(err)
	a.NotNil(db)
	a.NotNil(mock)
	defer db.Close()

	store, err := NewDefaultMySQLStore(&dbr.Connection{
		DB:            db,
		Dialect:       dialect.MySQL,
		EventReceiver: nil,
	})

	rows := sqlmock.NewRows([]string{"id", "value", "pub_date", "created_at", "updated_at"}).
		AddRow("LVL", 1.00, dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now())).
		AddRow("EUR", 2.00, dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now())).
		AddRow("USD", 3.00, dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()))

	mock.ExpectQuery("SELECT * FROM `currency` WHERE `pub_date` = (SELECT MAX(pub_date) FROM `currency`)").
		WillReturnRows(rows)

	cs, err := store.AllLatest(context.Background())
	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, 3)

	a.NoError(mock.ExpectationsWereMet())
}

func TestDefaultMySQLStore_HistoryByID(t *testing.T) {
	a := assert.New(t)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	a.NoError(err)
	a.NotNil(db)
	a.NotNil(mock)
	defer db.Close()

	store, err := NewDefaultMySQLStore(&dbr.Connection{
		DB:            db,
		Dialect:       dialect.MySQL,
		EventReceiver: nil,
	})

	rows := sqlmock.NewRows([]string{"id", "value", "pub_date", "created_at", "updated_at"}).
		AddRow("EUR", 1.00, dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now())).
		AddRow("EUR", 1.00, dbr.NewNullTime(time.Now().AddDate(0, 0, -1)), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now())).
		AddRow("EUR", 1.00, dbr.NewNullTime(time.Now().AddDate(0, 0, -2)), dbr.NewNullTime(time.Now()), dbr.NewNullTime(time.Now()))

	mock.ExpectQuery("SELECT * FROM `currency` WHERE id = 'EUR' ORDER BY pub_date DESC").
		WillReturnRows(rows)

	cs, err := store.HistoryByID(context.Background(), "EUR")
	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, 3)

	a.NoError(mock.ExpectationsWereMet())
}

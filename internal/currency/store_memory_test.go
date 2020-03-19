package currency_test

import (
	"context"
	"testing"
	"time"

	"github.com/agubarev/tetest/internal/currency"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMemoryStore_AllInOne(t *testing.T) {
	a := assert.New(t)

	s := currency.NewMemoryStore()
	a.NotNil(s)

	//---------------------------------------------------------------------------
	// creating test items
	//---------------------------------------------------------------------------
	cs, err := s.BulkCreate(context.Background(), []currency.Currency{
		{ID: "LVL", Value: 1.00, PubDate: dbr.NewNullTime(time.Now())},
		{ID: "EUR", Value: 2.00, PubDate: dbr.NewNullTime(time.Now())},
		{ID: "USD", Value: 3.00, PubDate: dbr.NewNullTime(time.Now())},
	})

	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, 3)

	//---------------------------------------------------------------------------
	// obtaining all latest stored values
	//---------------------------------------------------------------------------
	cs, err = s.AllLatest(context.Background())
	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, 3)

	//---------------------------------------------------------------------------
	// obtaining all stored values by ID
	//---------------------------------------------------------------------------
	cs, err = s.AllByID(context.Background(), "EUR")
	a.NoError(err)
	a.NotNil(cs)
	a.Len(cs, 1)
}

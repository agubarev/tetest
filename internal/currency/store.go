package currency

import (
	"context"
	"time"
)

// Store represents an API interface contract
// NOTE: this is a very simplified version, inteded
// for demonstration purposes only
type Store interface {
	BulkCreateDay(ctx context.Context, day Day) (items []Currency, err error)
	ListByID(ctx context.Context, id string) (items []Currency, err error)
	ListByDate(ctx context.Context, date time.Time) (items []Currency, err error)
}

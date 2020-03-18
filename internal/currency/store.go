package currency

import (
	"context"
	"time"
)

// Store represents an API interface contract
// NOTE: this is a very simplified version, inteded
// for demonstration purposes only
type Store interface {
	BulkCreate(ctx context.Context, cs []Currency) (_ []Currency, err error)
	ByDate(ctx context.Context, date time.Time) (cs []Currency, err error)
	ByID(ctx context.Context, id string) (cs []Currency, err error)
}

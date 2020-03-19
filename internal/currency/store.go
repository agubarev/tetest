package currency

import (
	"context"
)

// Store represents an API interface contract
// NOTE: this is a very simplified version, inteded
// for demonstration purposes only
type Store interface {
	BulkCreate(ctx context.Context, cs []Currency) (_ []Currency, err error)
	AllLatest(ctx context.Context) (cs []Currency, err error)
	AllByID(ctx context.Context, id string) (cs []Currency, err error)
}

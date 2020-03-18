package currency

import (
	"strings"

	"github.com/gocraft/dbr/v2"
)

// Currency represents a single currency item
type Currency struct {
	ID        string       `db:"id" json:"id"`
	Value     float64      `db:"value" json:"value"`
	ValidDate dbr.NullTime `db:"valid_date" json:"valid_date"`
	CreatedAt dbr.NullTime `db:"created_at" json:"created_at"`
}

func (c Currency) validate() (err error) {
	if strings.TrimSpace(c.ID) == "" {
		err = ErrEmptyCurrencyID
	}

	// NOTE: technically this could be zero, but very unlikely
	if c.Value == 0 {
		err = ErrInvalidCurrencyValue
	}

	return nil
}

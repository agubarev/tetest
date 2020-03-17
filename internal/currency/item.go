package currency

import (
	"strings"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// Day represents a simple container that groups
// currency values for a specific date
type Day struct {
	PublishedAt time.Time
	Roster      []Currency
}

func (day Day) validate() (err error) {
	for _, c := range day.Roster {
		if err = c.validate(); err != nil {
			return errors.Wrapf(err, "currency validation failed for date: %s", day.PublishedAt.Local().Format("02012006"))
		}
	}

	return nil
}

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

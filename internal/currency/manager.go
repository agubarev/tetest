package currency

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// errors
var (
	ErrNilDatabase          = errors.New("database is nil")
	ErrNilCurrencyStore     = errors.New("currency store is nil")
	ErrCurrencyNotFound     = errors.New("currency not found")
	ErrNoData               = errors.New("no data")
	ErrEmptyFeedURL         = errors.New("invalid feed url")
	ErrInvalidPayloadFormat = errors.New("invalid payload format")
	ErrEmptyCurrencyID      = errors.New("invalid currency id")
	ErrInvalidCurrencyValue = errors.New("invalid currency value")
)

// Manager handles business logic of its underlying objects
type Manager struct {
	feedAddr string
	store    Store
	logger   *zap.Logger
}

// NewCurrencyManager initializes a new manager
func NewManager(s Store, feedURL string) (*Manager, error) {
	if s == nil {
		return nil, errors.Wrap(ErrNilCurrencyStore, "failed to initialize currency manager")
	}

	feedURL = strings.TrimSpace(feedURL)
	if feedURL == "" {
		return nil, ErrEmptyFeedURL
	}

	m := &Manager{
		store:    s,
		feedAddr: feedURL,
	}

	return m, nil
}

// Store returns store if set
func (m *Manager) Store() (Store, error) {
	if m.store == nil {
		return nil, ErrNilCurrencyStore
	}

	return m.store, nil
}

// SetLogger assigns a primary logger for the manager
func (m *Manager) SetLogger(logger *zap.Logger) error {
	// if logger is set, then giving it a name
	// to know the log context
	if logger != nil {
		logger = logger.Named("[currency]")
	}

	m.logger = logger

	return nil
}

// Logger returns primary logger if is set, otherwise
// initializing and returning a new default emergency logger
// NOTE: will panic if it finally fails to obtain a logger
func (m *Manager) Logger() *zap.Logger {
	if m.logger == nil {
		l, err := zap.NewDevelopment()
		if err != nil {
			// having a working logger is crucial, thus must panic() if initialization fails
			panic(errors.Wrap(err, "failed to initialize fallback logger"))
		}

		m.logger = l
	}

	return m.logger
}

// Import imports external feed and returns as localized currency items
func (m *Manager) Import(ctx context.Context) (err error) {
	// initializing and parsing the feed
	f, err := gofeed.NewParser().ParseURL(m.feedAddr)
	if err != nil {
		log.Fatalf("failed to parse feed [%s]: %s", m.feedAddr, err)
	}

	//---------------------------------------------------------------------------
	// this function transforms raw currency payload into an intermediate map
	//---------------------------------------------------------------------------
	fn := func(s []string) (cs map[string]float64, err error) {
		slen := len(s)

		// slice must not be empty or contain an odd number of items
		if slen == 0 || slen%2 != 0 {
			return nil, ErrInvalidPayloadFormat
		}

		// initializing result map
		cs = make(map[string]float64, len(s)/2)

		// pairing values: key -> value
		// NOTE: returning if some value fails to be parsed
		for i := 0; i < slen; i += 2 {
			k, v := s[i], s[i+1]

			fval, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse value: %s", v)
			}

			cs[k] = fval
		}

		return cs, nil
	}

	//---------------------------------------------------------------------------
	// transforming raw payload into local currency
	//---------------------------------------------------------------------------
	for _, v := range f.Items {
		parsedMap, err := fn(strings.Split(strings.TrimSpace(v.Description), " "))
		if err != nil {
			return errors.Wrap(err, "failed to parse raw currency payload")
		}

		// initialzing currency slice
		cs := make([]Currency, 0, len(parsedMap))

		// initializing currency objects
		for id, value := range parsedMap {
			cs = append(cs, Currency{
				ID:      strings.ToUpper(id),
				Value:   value,
				PubDate: dbr.NewNullTime(v.PublishedParsed.Local()),
			})
		}

		// creating objects in bulk
		if _, err = m.BulkCreate(ctx, cs); err != nil {
			return errors.Wrapf(err, "failed to import currency for date: %s", v.PublishedParsed.Local().Format("02012006"))
		}
	}

	return nil
}

// BulkCreate creates currency values grouped by its publication date
func (m *Manager) BulkCreate(ctx context.Context, cs []Currency) (_ []Currency, err error) {
	// obtaining store
	store, err := m.Store()
	if err != nil {
		return nil, err
	}

	// validating and initializing new records
	for _, c := range cs {
		if err = c.validate(); err != nil {
			return nil, err
		}

		c.CreatedAt = dbr.NewNullTime(time.Now())
	}

	// storing items
	cs, err = store.BulkCreate(ctx, cs)
	if err != nil {
		return cs, err
	}

	return cs, nil
}

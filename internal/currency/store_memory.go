package currency

import (
	"context"
	"sync"
	"time"

	"github.com/gocraft/dbr/v2"
)

type defaultMemoryStore struct {
	cs map[string]map[string]Currency
	sync.RWMutex
}

func NewMemoryStore() Store {
	return &defaultMemoryStore{
		cs: make(map[string]map[string]Currency),
	}
}

func (s defaultMemoryStore) BulkCreate(ctx context.Context, cs []Currency) (_ []Currency, err error) {
	if s.cs == nil {
		panic("in-memory store is nil")
	}

	if len(cs) == 0 {
		return cs, nil
	}

	s.Lock()

	// adding given items to the runtime cache
	for k := range cs {
		c := &cs[k]

		pubDateKey := c.PubDate.Time.Format("2006-02-01")

		// initializing inner map if it hasn't been done yet
		if s.cs[pubDateKey] == nil {
			s.cs[pubDateKey] = make(map[string]Currency)
		}

		// assigning timestamp depending on whether this
		// currency is already in the store
		if _, ok := s.cs[pubDateKey][c.ID]; !ok {
			c.CreatedAt = dbr.NewNullTime(time.Now())
		} else {
			c.UpdatedAt = dbr.NewNullTime(time.Now())
		}

		// caching currency
		s.cs[pubDateKey][c.ID] = *c
	}

	s.Unlock()

	return cs, nil
}

func (s defaultMemoryStore) AllLatest(ctx context.Context) (cs []Currency, err error) {
	if s.cs == nil {
		panic("in-memory store is nil")
	}

	//---------------------------------------------------------------------------
	// finding latest stored publication date; the initial date is -100 years
	// (I hope I could debug in 100 years when it could become a problem) ^^
	//---------------------------------------------------------------------------
	latestPubDate := time.Now().AddDate(-100, 0, 0)

	for i := range s.cs {
		for j := range s.cs[i] {
			// first found currency's pubdate
			cpd := s.cs[i][j].PubDate.Time.Local()

			if latestPubDate.Before(cpd) {
				latestPubDate = cpd
			}
		}
	}

	// contains all currencies of the latest pub. date
	latestDay := s.cs[latestPubDate.Format("2006-02-01")]

	// initialzing result slice
	cs = make([]Currency, 0, len(latestDay))

	// copying values from the stored map into the result
	for k := range latestDay {
		cs = append(cs, latestDay[k])
	}

	return cs, nil
}

func (s defaultMemoryStore) AllByID(ctx context.Context, id string) (cs []Currency, err error) {
	if s.cs == nil {
		panic("in-memory store is nil")
	}

	// NOTE: not checking the validity of an ID

	// initialzing result
	cs = make([]Currency, 0)

	for pubDate := range s.cs {
		if c, ok := s.cs[pubDate][id]; ok {
			cs = append(cs, c)
		}
	}

	return cs, nil
}

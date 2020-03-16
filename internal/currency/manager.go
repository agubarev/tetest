package currency

import (
	"strings"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// errors
var (
	ErrEmptyFeedURL = errors.New("invalid feed url")
)

// Manager handles all currency related operations
type Manager struct {
	feedAddr string
}

// Item represents a single currency item
type Item struct {
	ID         string       `db:"id" json:"id"`
	Value      float64      `db:"value" json:"value"`
	RevisionID int          `db:"revision_id" json:"revision_id"`
	Date       dbr.NullTime `db:"date" json:"date"`
	CreatedAt  dbr.NullTime `db:"created_at" json:"created_at"`
}

// NewManager initializes a new currency manager
func NewManager(feedURL string) (*Manager, error) {
	feedURL = strings.TrimSpace(feedURL)
	if feedURL == "" {
		return nil, ErrEmptyFeedURL
	}

	return &Manager{feedAddr: feedURL}, nil
}

// Import imports external feed and returns
// as localized currency items
func (m *Manager) Import() ([]Item, error) {

}

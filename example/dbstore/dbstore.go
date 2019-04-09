package dbstore

import (
	"github.com/dimes/di/embeds"
)

// DBStore is an interface for Database interaction
type DBStore interface {
	StoreString(value string) error
	GetString() string
}

// TableDBStore is a DBStore implementation that writes to memory
type TableDBStore struct {
	inject embeds.Inject

	value string
}

// StoreString stores a string in the table
func (t *TableDBStore) StoreString(value string) error {
	t.value = value
	return nil
}

// GetString returns the in-memory string value
func (t *TableDBStore) GetString() string {
	return t.value
}

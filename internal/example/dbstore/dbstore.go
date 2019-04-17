package dbstore

import (
	"github.com/dimes/dihedral/embeds"
)

// DBStore is an interface for Database interaction
type DBStore interface {
	StoreString(value string) error
	GetString() string
}

// Prefix is a prefix to append to calls to GetString
type Prefix string

// MemoryDBStore is a DBStore implementation that writes to memory
type MemoryDBStore struct {
	inject embeds.Inject
	Prefix Prefix

	value string
}

// StoreString stores a string in the table
func (m *MemoryDBStore) StoreString(value string) error {
	m.value = value
	return nil
}

// GetString returns the in-memory string value
func (m *MemoryDBStore) GetString() string {
	return string(m.Prefix) + " " + m.value
}

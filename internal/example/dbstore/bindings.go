package dbstore

import (
	"github.com/dimes/dihedral/embeds"
)

// DBBindingModule binds the MemoryDBStore to the DBStore interface for injection
type DBBindingModule interface {
	BindsTableDBStore(impl *MemoryDBStore) DBStore

	// Interface modules can declare dependencies on other modules
	Modules() *DBProviderModule
}

// DBProviderModule provides DB configuration
type DBProviderModule struct {
	provide embeds.ProvidedModule
	Prefix  Prefix
}

// ProvidesPrefix provides the prefix to use
func (d *DBProviderModule) ProvidesPrefix() Prefix {
	return d.Prefix
}

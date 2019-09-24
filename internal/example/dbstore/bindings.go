package dbstore

import (
	"github.com/dimes/dihedral/embeds"
)

// DBProviderPrefix is a Prefix provided by the DBProviderModule
type DBProviderPrefix Prefix

// DBBindingModule binds the MemoryDBStore to the DBStore interface for injection
type DBBindingModule interface {
	BindsTableDBStore(impl *MemoryDBStore) DBStore

	BindsPrefix(impl DBProviderPrefix) Prefix

	// Interface modules can declare dependencies on other modules
	Modules() *DBProviderModule
}

// DBProviderModule provides DB configuration
type DBProviderModule struct {
	provide embeds.ProvidedModule
	Prefix  Prefix
}

// ProvidesPrefix provides the prefix to use
func (d *DBProviderModule) ProvidesPrefix() DBProviderPrefix {
	return DBProviderPrefix(d.Prefix)
}

package dbstore

// DBBindingModule binds the TableDBStore to the DBStore interface for injection
type DBBindingModule interface {
	BindsTableDBStore(impl *TableDBStore) DBStore

	// Interface modules can declare dependencies on other modules
	Modules() *DBProviderModule
}

// DBProviderModule provides DB configuration
type DBProviderModule struct{}

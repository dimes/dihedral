package example

import (
	"time"

	"github.com/dimes/dihedral/embeds"
	"github.com/dimes/dihedral/internal/example/dbstore"
)

// ServiceTimeout is the amount of time the service has to handle the operation
type ServiceTimeout time.Duration

// Service is the service struct we ultimately want to inject
type Service struct {
	inject embeds.Inject // Mark this struct as automatically injectable

	ServiceTimeout ServiceTimeout
	DBStore        dbstore.DBStore
}

// SetValueInDBStore sets a value from the DB store
func (s *Service) SetValueInDBStore(value string) error {
	return s.DBStore.StoreString(value)
}

// GetValueFromDBStore gets a value from the DB store
func (s *Service) GetValueFromDBStore() string {
	return s.DBStore.GetString()
}

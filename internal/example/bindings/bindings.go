//go:generate dihedral -definition ServiceDefinition

// Package bindings setups the the dependency bindings for a component
package bindings

import (
	"time"

	"github.com/dimes/dihedral/internal/example"
	"github.com/dimes/dihedral/internal/example/dbstore"
)

// DI has five concepts:
// 1. Provider Modules are structs whose methods return types that can be injected.
// 2. Binding Modules are interfaces whose methods are of the form:
//        func (m *Module) Binds(implementation *StructType) InterfaceType
//    These modules define a mapping between an interface and an implementation.
//    This allows injection of interface types backed by a concrete implementation.
// 3. Injectable Structs are structs with a non-exported member of type embeds.Inject.
//    These structs can automatically be constructed by DI without being provided by a provider
//    module.
// 4. Components: A component is an interface that defines the top-level types that can
//    be injected.
// 5. InjectionDefinition: Contains a Target component and a list of modules to include for doing the
//    injection

// ServiceDefinition defines the target and the modules to include
type ServiceDefinition interface {
	// The list of modules to include
	Modules() (*ServiceModule, dbstore.DBBindingModule)

	// An implementation of the interface will be automatically generated. The values
	// returned will be automatically instiated from their dependencies.
	Target() ServiceComponent
}

// ServiceComponent defines the top level component. The return type of
// the `Modules()` method is used as a list of modules to include for injection
// All other methods are considered types to generate injections for
type ServiceComponent interface {
	// The actual instance to return (fully injected). Errors during injection
	// will be returned in the error
	GetService() (*example.Service, error)

	// Non-interface / pointer types cannot return an error
	GetServiceTimeout() (example.ServiceTimeout, error)
}

// ServiceModule illustrates how each method on a struct module can provide
// an instance to be injected
type ServiceModule struct{}

// ProvidesServiceTimeout provides a time.Duration under the name ServiceTimeout
func (s *ServiceModule) ProvidesServiceTimeout() (example.ServiceTimeout, error) {
	return example.ServiceTimeout(5 * time.Second), nil
}

//go:generate dihedral -component ServiceComponent

// Package bindings setups the the dependency bindings for a component
package bindings

import (
	"time"

	"github.com/dimes/dihedral/example"
	"github.com/dimes/dihedral/example/dbstore"
)

// DI has four concepts:
// 1. Provider Modules are structs whose methods return types that can be injected.
// 2. Binding Modules are interfaces whose methods are of the form:
//        func (m *Module) Binds(implementation *StructType) InterfaceType
//    These modules define a mapping between an interface and an implementation.
//    This allows injection of interface types backed by a concrete implementation.
// 3. Injectable Structs are structs with a non-exported member of type embeds.Inject.
//    These structs can automatically be constructed by DI with being provided by a provider
//    module.
// 4. Components: A component is an interface that defines the top-level types that can
//    be injected. The component contains a list of modules to include for doing the
//    injection

// ServiceComponent defines the top level component. The return type of
// the `Modules()` method is used as a list of modules to include for injection
// All other methods are considered types to generate injections for
type ServiceComponent interface {
	// The list of modules to include
	Modules() (*ServiceModule, dbstore.DBBindingModule)

	// The actual instance to return (fully injected)
	GetService() *example.Service
}

// ServiceModule illustrates how each method on a struct module can provide
// an instance to be injected
type ServiceModule struct{}

// ProvidesServiceTimeout provides a time.Duration under the name ServiceTimeout
func (s *ServiceModule) ProvidesServiceTimeout() example.ServiceTimeout {
	return example.ServiceTimeout(5 * time.Second)
}

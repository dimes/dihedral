//go:generate di -component MyComponent

package test

import (
	"github.com/dimes/di/embeds"
	"github.com/dimes/di/test/othermod"
)

type Greeting string
type ServiceTimeout int

type MyComponent interface {
	Modules() (*Module, *SecondModule, BindingModule, *othermod.Module)
	Target() *MyTarget
}

type Module struct {
	provided embeds.ProvidedModule
	Greeting string
}

type SecondModule struct {
}

func (m *Module) ProvidesGreeting() Greeting {
	return Greeting(m.Greeting)
}

func (m *Module) ProvideServiceTimeout() ServiceTimeout {
	return ServiceTimeout(10)
}

type BindingModule interface {
	BindsMyImplementation(impl *MyImplementation) MyInterface
}

type MyTarget struct {
	inject embeds.Inject

	Greeting    Greeting
	Timeout     ServiceTimeout
	MyInterface MyInterface
	Additional  *othermod.AdditionalFunctionality
}

type MyInterface interface {
	TestMethod() string
}

type MyImplementation struct {
	inject embeds.Inject
}

func (m *MyImplementation) TestMethod() string {
	return "test"
}

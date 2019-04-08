package test

import (
	"github.com/dimes/di/embeds"
)

type Greeting string
type ServiceTimeout int

type MyComponent interface {
	Modules() (*Module, BindingModule)
	Target() *MyTarget
}

type Module struct {
}

func (m *Module) ProvidesGreeting() Greeting {
	return Greeting("Hello")
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

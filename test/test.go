package test

import (
	"github.com/dimes/di/embeds"
)

type Greeting string
type ServiceTimeout int

type MyComponent interface {
	Modules() *Module
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

type MyTarget struct {
	inject embeds.Inject

	Greeting Greeting
	Timeout  ServiceTimeout
}

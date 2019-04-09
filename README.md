# di

**di** (pronounced dee) is a compile-time injection framework for Go.

# Getting started

    > go get github.com/dimes/di

Create a type you want injected

    type MyFieldType string
    type InjectMe struct {
        inject  embeds.Inject
        MyField MyFieldType
    }

Create a module to provide non-injected dependencies

    type MyModule struct {}
    func (m *MyModule) ProvidesMyField() MyFieldType {
        return MyFieldType("Hello there")
    }

Create a component as the root of the dependency injection

    interface Component {
        Modules() *MyModule 
        InjectMePlease() *InjectMe
    }

Generate the bindings

    > di -component Component

Use the bindings

    func main() {
        component := di.NewComponent()
        injected := component.InjectMePlease()
        fmt.Println(string(injected.MyField)) # Prints "Hello there"
    }

## Differences from Wire

Wire, Google's injection framework, is another compile time framework for Go. Both frameworks are inspired
by Dagger. di differs from Wire in that di focuses on auto-injected components and self-contained
modules, whereas Wire focuses more on type registration via provider functions. di also leverages
struct receivers for better organization of runtime provided types. This makes di much more pleasurable
to work with.



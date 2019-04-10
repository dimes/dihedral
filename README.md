# dihedral

**dihedral** is a compile-time injection framework for Go.

# Getting started

    > go get github.com/dimes/dihedral

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

    > dihedral -component Component

Use the bindings

    func main() {
        component := digen.NewComponent()
        injected := component.InjectMePlease()
        fmt.Println(string(injected.MyField)) # Prints "Hello there"
    }

See the [example](example/) for a more detailed overview.

### Differences from Wire

Wire, Google's injection framework, is another compile time framework for Go. Both frameworks are inspired
by Dagger. **dihedral** differs from Wire in that **dihedral** focuses on auto-injected components and self-contained modules, whereas Wire focuses more on type registration via provider functions. **dihedral** also leverages struct receivers for better organization of runtime provided types. This makes di much more pleasurable to work with.

**dihedral**'s component structure also enables one to have multiple injected components that share modules. The type annotation system allows for auto-injected components, provided modules, and, in the future, sub-components and have a different scope than the parent component.

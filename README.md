# dihedral

**dihedral** is a compile-time injection framework for Go.

# Getting started

    > go get github.com/dimes/dihedral

Create a type you want injected

    type ServiceEndpoint string  // Name this string "ServiceEndpoint"
    type Service struct {
        inject  embeds.Inject    // Auto-inject this struct 
        Endpoint ServiceEndpoint // Inject a string with name "ServiceEndpoint"
    }

Create a module to provide non-injected dependencies

    // Each public method on this struct provides a type
    type ServiceModule struct {}
    func (s *ServiceModule) ProvidesServiceEndpoint() ServiceEndpoint {
        return ServiceEndpoint("http://hello.world")
    }

Create a component as the root of the dependency injection

    // A component tells dihedral which modules to use and the root of the DI graph
    interface ServiceComponent {
        Modules() *MyModule      // Tells dihedral which modules to include
        InjectService() *Service // Tells dihedral the root of the DI graph
    }

Generate the bindings

    > dihedral -component ServiceComponent

Use the bindings

    func main() {
        // dihedral generates the digen package
        component := digen.ServiceComponent()
        service := component.InjectService()
        fmt.Println(string(injected.Endpoint)) # Prints "http://hello.world"
    }

See the [example](example/) for a more detailed overview.

### Differences from Wire

Wire, Google's injection framework, is another compile time framework for Go. Both frameworks are inspired
by Dagger. **dihedral** differs from Wire in that **dihedral** focuses on auto-injected components and self-contained modules, whereas Wire focuses more on type registration via provider functions. **dihedral** also leverages struct receivers for better organization of runtime provided types. This makes di much more pleasurable to work with.

**dihedral**'s component structure also enables one to have multiple injected components that share modules. The type annotation system allows for auto-injected components, provided modules, and, in the future, sub-components and have a different scope than the parent component.

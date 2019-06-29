# Dihedral

[![Build Status](https://travis-ci.org/dimes/dihedral.svg?branch=master)](https://travis-ci.org/dimes/dihedral)

**Dihedral** is a compile-time injection framework for Go.

# Getting started

    > go get -u github.com/dimes/dihedral

Create a type you want injected:

    type ServiceEndpoint string  // Name this string "ServiceEndpoint"
    type Service struct {
        inject  embeds.Inject    // Auto-inject this struct 
        Endpoint ServiceEndpoint // Inject a string with name "ServiceEndpoint"
    }

Create a module to provide non-injected dependencies:

    // Each public method on this struct provides a type
    type ServiceModule struct {}
    func (s *ServiceModule) ProvidesServiceEndpoint() ServiceEndpoint {
        return ServiceEndpoint("http://hello.world")
    }

Create a component as the root of the dependency injection:

    type ServiceComponent interface {
        InjectService() *Service
    }

Create a definition with the target component and the module configuration:

    type ServiceDefinition interface {
        Modules() (*ServiceModule, dbstore.DBBindingModule)
        Target() ServiceComponent
    }

Generate the bindings

    > dihedral -definition ServiceDefinition

Use the bindings

    func main() {
        // dihedral generates the digen package
        component := digen.NewDihedralServiceComponent()
        service := component.InjectService()
        fmt.Println(string(injected.Endpoint)) # Prints "http://hello.world"
    }

## Further Reading

* [Documentation](https://dimes.github.io/dihedral/docs/)
* [Example](internal/example/)

### Differences from Wire

Wire, Google's injection framework, is another compile time framework for Go. Both frameworks are inspired
by Dagger. **Dihedral** differs from Wire in that **Dihedral** focuses on auto-injected components and self-contained modules, whereas Wire focuses more on type registration via provider functions. **Dihedral** also leverages struct receivers for better organization of runtime provided types. These features make **Dihedral** nicer to work with. 

**Dihedral**'s component structure also enables one to have multiple injected components that share modules. The type annotation system allows for auto-injected components, provided modules, and, in the future, sub-components and have a different scope than the parent component.

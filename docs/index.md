---
layout: default
title: Dihedral
description: Compile-time dependency injection for Go
---

# Getting started

    > go get -u github.com/dimes/dihedral

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

For more information, read the [docs](/dihedral/docs).

---
layout: sidebar
title: Dihedral
---

## Binding Modules

Binding modules are interfaces that "bind" implementations to an interface. Each method on the interface is treated as a type binding and is expected to have one parameter and one return type. The return type is an interface, and the parameter type is a struct that implements that interface. During injection, instances of that interface will be provided by the bound implementation. 

Binding modules also have a special `Modules()` method. This method takes no parameters and returns a list of other modules to use. This is useful for encapsulating logic in modules that can then be used as an entire unit.

```
type MyBindingModule interface {
    Modules() (*MyProviderModule, MyOtherBindingModule)
    BindsDatabase(impl *SQLDatabase) Database
}
```

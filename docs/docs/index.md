---
layout: sidebar
title: Dihedral
---

## Concepts 

**Dihedral** is powered by four simple concepts.

#### 1. Provider Modules

Provider modules are structs with associated "provider" methods, which is an exported method that returns an instance of some type. These methods are called directly at runtime when the injection is performed. Provider functions can take parameters that are injected.

```
type MyProviderModule struct {}
func (m *MyProviderModule) ProvidesSQLDatabase(tableName TableName) *SQLDatabase {
    return &SQLDatabase{tableName: tableName}
}
```

#### 2. Binding Modules

Binding modules are interfaces that "bind" implementations to an interface. Each method on the interface is treated as a type binding and is expected to have one parameter and one return type. The return type is an interface, and the parameter type is a struct that implements that interface. During injection, instances of that interface will be provided by the bound implementation. 

Binding modules also have a special `Modules()` method. This method takes no parameters and returns a list of other modules to use. This is useful for encapsulating logic in modules that can then be used as an entire unit.

```
type MyBindingModule interface {
    Modules() (*MyProviderModule, MyOtherBindingModule)
    BindsDatabase(impl *SQLDatabase) Database
}
```

#### 3. Struct Injection

**Dihedral** can automatically inject structs. This reduces a huge amount of boilerplate and keeps provider modules very small. To inject a struct, it must contain a non-exported field of type `embeds.Inject`. By default, all exported fields are injected. To skip a field, add the tag `di:"-"`.

```
type Service struct {
    inject       embeds.Inject
    ServiceDB    Database
    RequestCount int           `di:"-"`
}
```

#### 4. Components

Components tie everything together. Ultimately, **Dihedral** generates a component struct and constructor function. The methods on a component return fully injected types. Like binding modules, component interfaces also include a `Modules()` method that specifies which modules to use for injection.

```
type ServiceComponent interface {
    Modules() (*MyProviderModule, MyBindingModule)
    InjectService() *Service
}
```

---
layout: sidebar
title: Dihedral
---

## Components

Components defined what gets injected. Ultimately, **Dihedral** generates a component struct and constructor function that conforms to the component interface. The methods on a component return fully injected types. 

```
type ServiceComponent interface {
    InjectService() *Service
}
```

Once the code is generated, an implementation of `ServiceComponent` can be created using `digen.NewDihedralServiceComponent()`.

## Definitions

Definitions define the configuration for the injection. A definition is an interface that specifies
the target component for the injection. Like binding modules, definition interfaces also include a `Modules()` method that specifies which modules to use for injection.

```
type ServiceDefinition interface {
	Modules() (*ServiceModule, dbstore.DBBindingModule)
	Target() ServiceComponent
}
```

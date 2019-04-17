---
layout: sidebar
title: Dihedral
---

## 4. Components

Components tie everything together. Ultimately, **Dihedral** generates a component struct and constructor function. The methods on a component return fully injected types. Like binding modules, component interfaces also include a `Modules()` method that specifies which modules to use for injection.

```
type ServiceComponent interface {
    Modules() (*MyProviderModule, MyBindingModule)
    InjectService() *Service
}
```

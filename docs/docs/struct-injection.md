---
layout: sidebar
title: Dihedral
---

## Struct Injection

**Dihedral** can automatically inject structs. This reduces a huge amount of boilerplate and keeps provider modules very small. To inject a struct, it must contain a non-exported field of type `embeds.Inject`. By default, all exported fields are injected. To skip a field, add the tag `di:"-"`.

```
type Service struct {
    inject       embeds.Inject
    ServiceDB    Database
    RequestCount int           `di:"-"`
}
```
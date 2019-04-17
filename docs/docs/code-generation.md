---
layout: sidebar
title: Dihedral
---

## Code Generation

Code can be generated using the CLI:

    > dihedral -component ServiceComponent

This can also be used with `go generate` by putting the following at the top of the file containing the component:

    //go:generate dihedral -component ServiceComponent

The generated code is located in the `digen` folder. If the component is named `ServiceComponent`, then **dihedral** generates a function named `func NewServiceComponent() *GeneratedServiceComponent` that contains the requested injections. 

```
func main() {
    component := digen.NewServiceComponent()
    service := component.InjectService()
}
```

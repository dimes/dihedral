---
layout: sidebar
title: Dihedral
---

## Code Generation

Code can be generated using the CLI:

    > dihedral -definition ServiceDefinition

This can also be used with `go generate` by putting the following at the top of the file containing the component:

    //go:generate dihedral -definition ServiceDefinition

The generated code is located in the `digen` folder. If the definition has a target component named `ServiceComponent`, then **dihedral** generates a function named `func NewDihedralServiceComponent() *DihedralServiceComponent` that contains the requested injections. 

```
func main() {
    component := digen.NewDihedralServiceComponent()
    service := component.InjectService()
}
```

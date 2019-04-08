// Package structs contains common type definitions
package structs

import "go/types"

// Interface contains a *type.Interface as well as a qualifier
type Interface struct {
	Name *types.Named
	Type *types.Interface
}

// Struct contains a *type.Struct as well as a qualifier
type Struct struct {
	Name *types.Named
	Type *types.Struct
}

package gen

import (
	"fmt"
	"go/types"

	"github.com/dimes/di/resolver"
	"github.com/dimes/di/structs"
	"github.com/dimes/di/typeutil"
)

// FactoryName returns the name of the factory function for the given name
func FactoryName(typeName *types.Named) string {
	return "factory_" + SanitizeName(typeName)
}

// ProviderName returns the name of the provider function for the given name
func ProviderName(typeName *types.Named) string {
	return "provides_" + SanitizeName(typeName)
}

// Assignment represents a way of getting a injected value, either by a provider
// or by an injectable factory method
type Assignment interface {
	// ID returns the ID of the type being assigned
	ID() string

	// GetSourceAssignment returns the assignment as a string of source code
	GetSourceAssignment() string
}

type factoryAssignment struct {
	typeName *types.Named
}

// NewFactoryAssignment returns a factory-method based assignment
func NewFactoryAssignment(typeName *types.Named) Assignment {
	return &factoryAssignment{
		typeName: typeName,
	}
}

func (f *factoryAssignment) ID() string {
	return typeutil.IDFromNamed(f.typeName)
}

func (f *factoryAssignment) GetSourceAssignment() string {
	return FactoryName(f.typeName) + "(" + componentName + ")"
}

type providerAssignment struct {
	typeName *types.Named
}

// NewProviderAssignment returns a component provided assignment
func NewProviderAssignment(typeName *types.Named) Assignment {
	return &providerAssignment{
		typeName: typeName,
	}
}

func (p *providerAssignment) ID() string {
	return typeutil.IDFromNamed(p.typeName)
}

func (p *providerAssignment) GetSourceAssignment() string {
	return componentName + "." + ProviderName(p.typeName) + "()"
}

// AssignmentForFieldType returns an assignment for the given field type
func AssignmentForFieldType(
	rawFieldType types.Type,
	providers map[string]*resolver.ResolvedType,
	bindings map[string]*structs.Struct,
) (Assignment, error) {
	var fieldName *types.Named
	switch fieldType := rawFieldType.(type) {
	case *types.Named:
		fieldName = fieldType
	case *types.Pointer:
		fieldName = fieldType.Elem().(*types.Named)
	default:
		return nil, fmt.Errorf("Field %+v is not a supported type", fieldType)
	}

	fieldID := typeutil.IDFromNamed(fieldName)
	if binding := bindings[fieldID]; binding != nil {
		fieldID = typeutil.IDFromNamed(binding.Name)
	}

	if provider := providers[fieldID]; provider != nil {
		return NewProviderAssignment(fieldName), nil
	}

	return NewFactoryAssignment(fieldName), nil
}

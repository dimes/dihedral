package gen

import (
	"fmt"
	"go/types"

	"github.com/dimes/dihedral/resolver"
	"github.com/dimes/dihedral/typeutil"
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
	// CastTo returns the type this assignment should be cast to. If nil is returned, then
	// the assignment is not cast to any type.
	CastTo() *types.Named

	// GetSourceAssignment returns the assignment as a string of source code.
	GetSourceAssignment() string
}

type factoryAssignment struct {
	componentReceiverName string
	typeName              *types.Named
}

// NewFactoryAssignment returns a factory-method based assignment
func NewFactoryAssignment(
	componentReceiverName string,
	typeName *types.Named,
) Assignment {
	return &factoryAssignment{
		componentReceiverName: componentReceiverName,
		typeName:              typeName,
	}
}

func (f *factoryAssignment) CastTo() *types.Named {
	return nil
}

func (f *factoryAssignment) GetSourceAssignment() string {
	return FactoryName(f.typeName) + "(" + f.componentReceiverName + ")"
}

type providerAssignment struct {
	componentReceiverName string
	typeName              *types.Named
	castTo                *types.Named
}

// NewProviderAssignment returns a component provided assignment
func NewProviderAssignment(
	componentReceiverName string,
	typeName *types.Named,
	castTo *types.Named,
) Assignment {
	return &providerAssignment{
		componentReceiverName: componentReceiverName,
		typeName:              typeName,
		castTo:                castTo,
	}
}

func (p *providerAssignment) CastTo() *types.Named {
	return p.castTo
}

func (p *providerAssignment) GetSourceAssignment() string {
	return p.componentReceiverName + "." + ProviderName(p.typeName) + "()"
}

// AssignmentForFieldType returns an assignment for the given field type
func AssignmentForFieldType(
	componentReceiverName string,
	rawFieldType types.Type,
	providers map[string]resolver.ResolvedType,
	bindings map[string]*types.Named,
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

	var castTo *types.Named
	fieldID := typeutil.IDFromNamed(fieldName)
	if binding := bindings[fieldID]; binding != nil {
		if fieldName != binding {
			castTo = fieldName
		}

		fieldID = typeutil.IDFromNamed(binding)
		fieldName = binding
	}

	if provider := providers[fieldID]; provider != nil {
		typedProvider, ok := provider.(*resolver.ModuleResolvedType)
		if ok {
			fmt.Printf("%+v \t %+v \t %+v\n", rawFieldType, fieldName, castTo)
			fieldName = typedProvider.Name
			return NewProviderAssignment(componentReceiverName, fieldName, castTo), nil
		}

		return nil, fmt.Errorf("Unknown provider type %+v", provider)
	}

	return NewFactoryAssignment(componentReceiverName, fieldName), nil
}

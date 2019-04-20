package gen

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/dimes/dihedral/embeds"
	"github.com/dimes/dihedral/resolver"
	"github.com/dimes/dihedral/structs"
	"github.com/dimes/dihedral/typeutil"
	"github.com/pkg/errors"
)

const (
	diTag   = "di"
	skipTag = "-"
)

var (
	injectType = reflect.TypeOf(embeds.Inject{})
)

// GeneratedFactory contains information for generating a factory for
// an injected struct
//
// The generated code from this factory looks something like this:
//
// func TargetFactory(component *GeneratedComponent) *TargetType {
//     target := &TargetType{}
//     targetType.ProvidedType = component.provides_ProvidedType()
//     targetType.InjectableType = InjectableFactory(component)
//     return target
// }
type GeneratedFactory struct {
	targetName   *types.Named
	targetStruct *types.Struct
	assignments  map[string]Assignment
	dependencies []*injectionTarget
}

// NewGeneratedFactoryIfNeeded generates a factory for the given struct.
// If a factory cannot be generated, e.g. if the struct is not injectable,
// nil is returned
func NewGeneratedFactoryIfNeeded(
	targetName *types.Named,
	targetStruct *types.Struct,
	providers map[string]resolver.ResolvedType,
	bindings map[string]*structs.Struct,
) (*GeneratedFactory, error) {
	if targetStruct == nil {
		return nil, nil
	}

	if !typeutil.HasFieldOfType(targetStruct, injectType) {
		return nil, nil
	}

	assignments := make(map[string]Assignment)
	dependencies := make([]*injectionTarget, 0)
	for i := 0; i < targetStruct.NumFields(); i++ {
		field := targetStruct.Field(i)
		if !field.Exported() {
			continue
		}

		tags := strings.Split(reflect.StructTag(targetStruct.Tag(i)).Get(diTag), ",")
		if len(tags) > 0 && tags[0] == skipTag {
			continue
		}

		assignment, err := AssignmentForFieldType(field.Type(), providers, bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error generating bindings for %+v", targetStruct)
		}

		assignments[field.Name()] = assignment
		dependencies = append(dependencies, newInjectionTarget(field.Type()))
	}

	return &GeneratedFactory{
		targetName:   targetName,
		targetStruct: targetStruct,
		assignments:  assignments,
		dependencies: dependencies,
	}, nil
}

// ToSource converts this generated factory into Go source code. The
// source should be treated as a separate source file in the generated
// component package
func (g *GeneratedFactory) ToSource(componentPackage string) string {
	var builder strings.Builder
	builder.WriteString("package " + componentPackage + "\n")
	builder.WriteString("import target_pkg \"" + g.targetName.Obj().Pkg().Path() + "\"\n")
	builder.WriteString(
		"func " + FactoryName(g.targetName) + "(" + componentName + " *" + componentType +
			") *target_pkg." + g.targetName.Obj().Name() + " {\n")
	builder.WriteString("\ttarget := &target_pkg." + g.targetName.Obj().Name() + "{}\n")

	for name, assignment := range g.assignments {
		builder.WriteString("\ttarget." + name + " = " + assignment.GetSourceAssignment() + "\n")
	}

	builder.WriteString("\treturn target\n")
	builder.WriteString("}\n")

	return builder.String()
}

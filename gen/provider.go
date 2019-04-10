package gen

import (
	"go/types"
	"strings"

	"github.com/dimes/dihedral/resolver"
	"github.com/dimes/dihedral/structs"
	"github.com/pkg/errors"
)

// GeneratedProvider is a single generated provider method on the component
type GeneratedProvider struct {
	resolvedType *resolver.ResolvedType
	assignments  []Assignment
	dependencies []*injectionTarget
}

// NewGeneratedProvider generates a provider function for the given resolved type
// The generated function has the form:
//
// func (generatedComponent *GeneratedComponent) provides_Name() *SomeType {
//     return someModule.providerFunc(
//	       component.provides_ProvidedType(),
//         InjectableFactory(component),
//     )
// }
func NewGeneratedProvider(
	resolvedType *resolver.ResolvedType,
	providers map[string]*resolver.ResolvedType,
	bindings map[string]*structs.Struct,
) (*GeneratedProvider, error) {
	assignments := make([]Assignment, 0)
	dependencies := make([]*injectionTarget, 0)
	signature := resolvedType.Method.Type().(*types.Signature)
	for i := 0; i < signature.Params().Len(); i++ {
		param := signature.Params().At(i)
		assignment, err := AssignmentForFieldType(param.Type(), providers, bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error generating binding for %+v", resolvedType)
		}

		assignments = append(assignments, assignment)
		dependencies = append(dependencies, newInjectionTarget(param.Type()))
	}

	return &GeneratedProvider{
		resolvedType: resolvedType,
		assignments:  assignments,
		dependencies: dependencies,
	}, nil
}

// ToSource returns the source code for this provider.
func (g *GeneratedProvider) ToSource(componentPackage string) string {
	moduleVariableName := SanitizeName(g.resolvedType.Module.Name)
	returnType := "target_pkg." + g.resolvedType.Name.Obj().Name()
	if g.resolvedType.IsPointer {
		returnType = "*" + returnType
	}

	var builder strings.Builder
	builder.WriteString("package " + componentPackage + "\n")
	builder.WriteString("import target_pkg \"" + g.resolvedType.Name.Obj().Pkg().Path() + "\"\n")
	builder.WriteString(
		"func (" + componentName + " *" + componentType + ") " + ProviderName(g.resolvedType.Name) + "() " +
			returnType + " {\n")
	builder.WriteString(
		"\treturn " + componentName + "." + moduleVariableName + "." + g.resolvedType.Method.Name() + "(\n")

	for _, assignment := range g.assignments {
		builder.WriteString("\t\t" + assignment.GetSourceAssignment() + ",\n")
	}
	builder.WriteString("\t)\n")

	builder.WriteString("}\n")
	return builder.String()
}

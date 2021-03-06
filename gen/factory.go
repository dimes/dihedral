package gen

import (
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/dimes/dihedral/embeds"
	"github.com/dimes/dihedral/resolver"
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
	generatedComponentType     string
	generatedComponentReceiver string
	targetName                 *types.Named
	targetStruct               *types.Struct
	assignments                map[string]Assignment
	dependencies               []*injectionTarget
}

// NewGeneratedFactoryIfNeeded generates a factory for the given struct.
// If a factory cannot be generated, e.g. if the struct is not injectable,
// nil is returned
func NewGeneratedFactoryIfNeeded(
	generatedComponentType string,
	generatedComponentReceiver string,
	targetName *types.Named,
	targetStruct *types.Struct,
	providers map[string]resolver.ResolvedType,
	bindings map[string]*types.Named,
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

		assignment, err := AssignmentForFieldType(
			generatedComponentReceiver,
			field.Type(),
			providers,
			bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error generating bindings for %+v", targetStruct)
		}

		assignments[field.Name()] = assignment
		dependencies = append(dependencies, newInjectionTarget(field.Type()))
	}

	return &GeneratedFactory{
		generatedComponentType:     generatedComponentType,
		generatedComponentReceiver: generatedComponentReceiver,
		targetName:                 targetName,
		targetStruct:               targetStruct,
		assignments:                assignments,
		dependencies:               dependencies,
	}, nil
}

// ToSource converts this generated factory into Go source code. The
// source should be treated as a separate source file in the generated
// component package
func (g *GeneratedFactory) ToSource(componentPackage string) string {
	returnType := "target_pkg." + g.targetName.Obj().Name()
	var builder strings.Builder
	builder.WriteString("// Code generated by go generate; DO NOT EDIT.\n")
	builder.WriteString("package " + componentPackage + "\n")

	imports := map[string]string{
		g.targetName.Obj().Pkg().Path(): "target_pkg",
	}

	for _, assignment := range g.assignments {
		castTo := assignment.CastTo()
		if castTo == nil {
			continue
		}

		packagePath := castTo.Obj().Pkg().Path()
		if importName := imports[packagePath]; importName == "" {
			imports[packagePath] = "di_import_" + strconv.Itoa(len(imports)+1)
		}
	}

	builder.WriteString("import (\n")
	for packagePath, importName := range imports {
		builder.WriteString("\t" + importName + " \"" + packagePath + "\"\n")
	}
	builder.WriteString(")\n")

	builder.WriteString(
		"func " + FactoryName(g.targetName) +
			"(" + g.generatedComponentReceiver + " *" + g.generatedComponentType +
			") (*" + returnType + ", error) {\n")
	builder.WriteString("\ttarget := &" + returnType + "{}\n")

	paramCounter := 0
	for name, assignment := range g.assignments {
		paramName := fmt.Sprintf("param%d", paramCounter)
		paramCounter++

		builder.WriteString("\t" + paramName + ", err := " + assignment.GetSourceAssignment() + "\n")
		builder.WriteString("\tif err != nil {\n")
		builder.WriteString("\t\tvar zeroValue *" + returnType + "\n")
		builder.WriteString("\t\treturn zeroValue, err\n")
		builder.WriteString("\t}\n")

		sourceAssignment := paramName
		castTo := assignment.CastTo()
		if castTo != nil {
			importName := imports[castTo.Obj().Pkg().Path()]
			if importName != "" {
				importName = importName + "."
			}
			sourceAssignment = "(" + importName + castTo.Obj().Name() + ")(" + sourceAssignment + ")"
		}

		builder.WriteString("\ttarget." + name + " = " + sourceAssignment + "\n")
	}

	builder.WriteString("\treturn target, nil\n")
	builder.WriteString("}\n")

	return builder.String()
}

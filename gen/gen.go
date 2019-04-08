// Package gen contains the logic for generating the source code
package gen

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"

	"github.com/dimes/di/resolver"
	"github.com/dimes/di/structs"
	"github.com/dimes/di/typeutil"
	"github.com/pkg/errors"
)

const (
	// All generated references to the component must refere to it by this name
	componentName = "generatedComponent"

	// The type of the generated component
	componentType = "GeneratedComponent"
)

// GeneratedComponent is the resolve of GenerateComponent and contains helper methods
// for converting this directly into source
type GeneratedComponent struct {
	name                  string
	targetsAndAssignments []*targetAndAssignment
	factories             []*GeneratedFactory
	providers             []*GeneratedProvider
}

type injectionTarget struct {
	Type types.Type
}

type targetAndAssignment struct {
	target     *resolver.InjectionTarget
	assignment Assignment
}

func newInjectionTarget(targetType types.Type) *injectionTarget {
	return &injectionTarget{
		Type: targetType,
	}
}

// NewGeneratedComponent generates the source for the given component
func NewGeneratedComponent(
	componentName string,
	targets []*resolver.InjectionTarget,
	providers map[string]*resolver.ResolvedType,
	bindings map[string]*structs.Struct,
) (*GeneratedComponent, error) {
	seenTargets := make(map[string]struct{})

	injectionStack := make([]*injectionTarget, 0)
	for _, target := range targets {
		injectionStack = append(injectionStack, newInjectionTarget(target.Type))
	}

	factories := make([]*GeneratedFactory, 0)
	providerFuncs := make([]*GeneratedProvider, 0)
	for len(injectionStack) > 0 {
		target := injectionStack[len(injectionStack)-1]
		injectionStack = injectionStack[:len(injectionStack)-1]

		var targetName *types.Named
		var targetStruct *types.Struct
		switch typedTarget := target.Type.(type) {
		case *types.Named:
			targetID := typeutil.IDFromNamed(typedTarget)
			if provider := providers[targetID]; provider != nil {
				targetName = typedTarget
				// No target struct for providers
			} else if boundType := bindings[targetID]; boundType != nil {
				targetName = boundType.Name
				targetStruct = boundType.Type
			} else {
				return nil, fmt.Errorf("No type binding found for %+v", target)
			}
		case *types.Pointer:
			targetName = typedTarget.Elem().(*types.Named)
			targetStruct = targetName.Underlying().(*types.Struct)
		default:
			return nil, fmt.Errorf("Target %+v is of an unsupported type", target)
		}

		targetID := typeutil.IDFromNamed(targetName)
		if _, ok := seenTargets[targetID]; ok {
			continue
		}
		seenTargets[targetID] = struct{}{}

		factory, err := NewGeneratedFactoryIfNeeded(targetName, targetStruct, providers, bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting factory for target %+v", targetStruct)
		}

		if factory != nil {
			factories = append(factories, factory)
			injectionStack = append(injectionStack, factory.dependencies...)
			continue
		}

		provider := providers[targetID]
		if provider == nil {
			return nil, fmt.Errorf("Target %+v is not marked as injectable and has no provider", target)
		}

		providerFunc, err := NewGeneratedProvider(provider, providers, bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting provider for %+v", provider)
		}

		providerFuncs = append(providerFuncs, providerFunc)
		injectionStack = append(injectionStack, providerFunc.dependencies...)
	}

	targetsAndAssignments := make([]*targetAndAssignment, 0)
	for _, target := range targets {
		assignment, err := AssignmentForFieldType(
			target.Type,
			providers,
			bindings)
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting toplevel target for %+v", target)
		}

		targetsAndAssignments = append(targetsAndAssignments, &targetAndAssignment{
			target:     target,
			assignment: assignment,
		})
	}

	return &GeneratedComponent{
		name:                  componentName,
		targetsAndAssignments: targetsAndAssignments,
		factories:             factories,
		providers:             providerFuncs,
	}, nil
}

// ToSource returns a map of file names to file contents that represent the generated
// source of this component
func (g *GeneratedComponent) ToSource(componentPackage string) map[string]string {
	imports := make(map[string]string)
	moduleStructParams := make([]*structs.Struct, 0)
	for _, provider := range g.providers {
		packagePath := provider.resolvedType.Module.Name.Obj().Pkg().Path()
		if _, ok := imports[packagePath]; ok {
			continue
		}

		moduleStructParams = append(moduleStructParams, provider.resolvedType.Module)
		imports[packagePath] = "di_import_" + strconv.Itoa(len(imports)+1)
	}

	for _, targetAssignment := range g.targetsAndAssignments {
		target := targetAssignment.target
		packagePath := target.Name.Obj().Pkg().Path()
		if _, ok := imports[packagePath]; ok {
			continue
		}

		imports[packagePath] = "di_import_" + strconv.Itoa(len(imports)+1)
	}

	var builder strings.Builder
	builder.WriteString("package " + componentPackage + "\n")

	builder.WriteString("import (\n")
	for packagePath, importName := range imports {
		builder.WriteString("\t" + importName + " \"" + packagePath + "\"\n")
	}

	builder.WriteString(")\n")

	builder.WriteString("type " + componentType + " struct {\n")
	for _, module := range moduleStructParams {
		moduleImportName := imports[module.Name.Obj().Pkg().Path()]
		moduleTypeName := module.Name.Obj().Name()
		moduleVariableName := SanitizeName(module.Name)
		builder.WriteString(
			"\t" + moduleVariableName + " *" + moduleImportName + "." + moduleTypeName + "\n")
	}
	builder.WriteString("}\n")

	builder.WriteString("func New" + g.name + "() *" + componentType + " {\n")
	builder.WriteString("\t return &" + componentType + "{\n")
	for _, module := range moduleStructParams {
		moduleImportName := imports[module.Name.Obj().Pkg().Path()]
		moduleTypeName := module.Name.Obj().Name()
		moduleVariableName := SanitizeName(module.Name)
		builder.WriteString(
			"\t\t" + moduleVariableName + ": &" + moduleImportName + "." + moduleTypeName + "{},\n")
	}
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	for _, targetAssignment := range g.targetsAndAssignments {
		target := targetAssignment.target
		importName := imports[target.Name.Obj().Pkg().Path()]
		targetTypeName := target.Name.Obj().Name()
		returnType := importName + "." + targetTypeName
		if target.IsPointer {
			returnType = "*" + returnType
		}
		assignment := targetAssignment.assignment
		builder.WriteString(
			"func (" + componentName + " *" + componentType + ") " + target.MethodName + "() " +
				returnType + " {\n")
		builder.WriteString("\treturn " + assignment.GetSourceAssignment() + "\n")
		builder.WriteString("}\n")
	}

	output := map[string]string{
		"component": builder.String(),
	}

	for _, factory := range g.factories {
		output[SanitizeName(factory.targetName)+"_Factory"] = factory.ToSource(componentPackage)
	}

	for _, provider := range g.providers {
		output[SanitizeName(provider.resolvedType.Name)+"_Provider"] = provider.ToSource(componentPackage)
	}

	return output
}

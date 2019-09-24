// Package resolver handles resolving dependencies
package resolver

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/dimes/dihedral/structs"
	"github.com/dimes/dihedral/typeutil"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

const (
	modulesFunc = "Modules"
	targetFunc  = "Target"
)

var (
	reservedMethods = map[string]struct{}{
		modulesFunc: struct{}{},
	}
)

type resolutionNode struct {
	parent   *resolutionNode
	nodeType types.Type
}

// InjectionTarget represents something that should be injected
type InjectionTarget struct {
	MethodName string
	Type       types.Type
	Name       *types.Named
	IsPointer  bool
	HasError   bool
}

// ResolvedType is an interface that represents a type provided by
// an injection source. Currently, the only injection source is
// via a provider module.
type ResolvedType interface {
	DebugInfo() string
}

// ModuleResolvedType represents a type that has been resolved via a module.
type ModuleResolvedType struct {
	Module    *structs.Struct
	Method    *types.Func
	Name      *types.Named
	IsPointer bool
	HasError  bool
}

// DebugInfo implements ResolvedType DebugInfo
func (m *ModuleResolvedType) DebugInfo() string {
	return fmt.Sprintf("Module: %+v, method: %+v, type name: %+v, isPointer: %t",
		m.Module, m.Method, m.Name, m.IsPointer)
}

// ResolveResult is the result of ResolveComponentModules
type ResolveResult struct {
	TargetInterfaceName string                  // Name of the Target interface
	Targets             []*InjectionTarget      // List of injection targets
	Providers           map[string]ResolvedType // Map of type to the provider of that type
	Bindings            map[string]*types.Named // Map of interface to concrete type
}

// ResolveComponentModules resolves the modules for the component interface.
// The return types are:
// - List of struct modules (used to provide concrete types)
// - List of interface modules (used to bind interfaces to implementations)
func ResolveComponentModules(
	fileSet *token.FileSet,
	componentInterface *structs.Interface,
) (
	*ResolveResult,
	error,
) {
	targetInterfaceName, targets, err := getTargetsFromInterface(componentInterface.Type)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting targets for %+v", componentInterface)
	}

	stack, err := getNodesFromInterface(componentInterface.Type, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting modules for %+v", componentInterface)
	}

	seen := make(map[string]struct{})
	providers := make(map[string]ResolvedType)
	bindings := make(map[string]*types.Named)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch typedNode := node.nodeType.(type) {
		case *types.Named:
			nodeInterface, ok := typedNode.Underlying().(*types.Interface)
			if !ok {
				return nil, fmt.Errorf("Expected node %+v to be pointer or interface", typedNode)
			}

			id := typeutil.IDFromNamed(typedNode)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}

			nodeModules, err := getNodesFromInterface(nodeInterface, node)
			if err != nil {
				return nil, errors.Wrapf(err, "Error getting dependencies for %+v", nodeInterface)
			}

			stack = append(stack, nodeModules...)

			bindingInterface := &structs.Interface{
				Name: typedNode,
				Type: nodeInterface,
			}

			moduleBindings, err := extractBindings(bindingInterface)
			if err != nil {
				return nil, errors.Wrapf(err, "Error extracting bindings in %+v", nodeInterface)
			}

			for id, boundStruct := range moduleBindings {
				if _, ok := bindings[id]; ok {
					return nil, fmt.Errorf("Binding %+v seen twice", id)
				}

				if _, ok := providers[id]; ok {
					return nil, fmt.Errorf("Binding %+v seen twice", id)
				}

				bindings[id] = boundStruct
			}
		case *types.Pointer:
			namedNode, ok := typedNode.Elem().(*types.Named)
			if !ok {
				return nil, fmt.Errorf("Expected pointer %+v to point to named element", typedNode)
			}

			id := typeutil.IDFromNamed(namedNode)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}

			structNode, ok := namedNode.Underlying().(*types.Struct)
			if !ok {
				return nil, fmt.Errorf("Expected pointer %+v to point to a struct", typedNode)
			}

			module := &structs.Struct{
				Name: namedNode,
				Type: structNode,
			}

			config := &packages.Config{
				Mode: packages.LoadSyntax,
			}

			pkgs, err := packages.Load(config, namedNode.Obj().Pkg().Path())
			if err != nil {
				return nil, errors.Wrapf(err, "Error loading %+v", namedNode)
			}

			for _, astPkg := range pkgs {
				for identifier, definition := range astPkg.TypesInfo.Defs {
					if !identifier.IsExported() {
						continue
					}

					funcDefinition, ok := definition.(*types.Func)
					if !ok {
						continue
					}

					signature := funcDefinition.Type().(*types.Signature)
					receiver := signature.Recv()
					if receiver == nil {
						continue
					}

					pointerReceiver, ok := receiver.Type().(*types.Pointer)
					if !ok {
						continue
					}

					receiverName := pointerReceiver.Elem().(*types.Named)
					receiverID := typeutil.IDFromNamed(receiverName)

					if receiverID != id {
						continue
					}

					if signature.Results().Len() == 0 || signature.Results().Len() > 2 {
						return nil, fmt.Errorf("Expected at most two results from %+v", signature)
					}

					hasError := false
					if signature.Results().Len() == 2 {
						errType, ok := signature.Results().At(1).Type().(*types.Named)
						if !ok {
							return nil, fmt.Errorf("Expected %+v to return an error", signature)
						}

						if errType.Obj().Pkg() != nil {
							return nil, fmt.Errorf("Expected %+v to return an error", signature)
						}

						if errType.Obj().Name() != "error" {
							return nil, fmt.Errorf("Expected %+v to return an error", signature)
						}

						hasError = true
					}

					result := signature.Results().At(0)
					var resultName *types.Named

					isPointer := false
					switch resultType := result.Type().(type) {
					case *types.Pointer:
						isPointer = true
						resultName = resultType.Elem().(*types.Named)
					case *types.Named:
						resultName = resultType
					default:
						return nil, fmt.Errorf("Result %+v is an unsupported type", result)
					}

					resultID := typeutil.IDFromNamed(resultName)
					if _, ok := bindings[resultID]; ok {
						return nil, fmt.Errorf("Binding %+v seen twice", resultID)
					}

					if _, ok := providers[resultID]; ok {
						return nil, fmt.Errorf("Binding %+v seen twice", resultID)
					}

					resolvedType := &ModuleResolvedType{
						Module:    module,
						Method:    funcDefinition,
						Name:      resultName,
						IsPointer: isPointer,
						HasError:  hasError,
					}

					providers[resultID] = resolvedType
				}
			}
		default:
			return nil, fmt.Errorf("%+v is not a recognized module type", typedNode)
		}
	}

	return &ResolveResult{
		TargetInterfaceName: targetInterfaceName,
		Targets:             targets,
		Providers:           providers,
		Bindings:            bindings,
	}, nil
}

func getTargetsFromInterface(
	interfaceType *types.Interface,
) (
	string,
	[]*InjectionTarget,
	error,
) {
	targetMethod := typeutil.GetInterfaceMethod(interfaceType, targetFunc)
	if targetMethod == nil {
		return "", nil, fmt.Errorf("%+v has no Target() method", interfaceType)
	}

	targetSignature := targetMethod.Type().(*types.Signature)
	if targetSignature.Params().Len() > 0 {
		return "", nil, fmt.Errorf("Target method %+v has arguments. Expected exactly 0", targetMethod)
	}

	if targetSignature.Results().Len() != 1 {
		return "", nil, fmt.Errorf("Expected exactly on return type on %+v", targetMethod)
	}

	targetNamedType, ok := targetSignature.Results().At(0).Type().(*types.Named)
	if !ok {
		return "", nil, fmt.Errorf("Return type of %+v is not a named type", targetSignature)
	}

	targetInterface, ok := targetNamedType.Underlying().(*types.Interface)
	if !ok {
		return "", nil, fmt.Errorf("Return type of %+v is not an interface", targetSignature)
	}

	targets := make([]*InjectionTarget, 0)
	for i := 0; i < targetInterface.NumMethods(); i++ {
		method := targetInterface.Method(i)
		if !method.Exported() {
			continue
		}

		if _, ok := reservedMethods[method.Name()]; ok {
			continue
		}

		signature := method.Type().(*types.Signature)
		if signature.Params().Len() > 0 {
			return "", nil, fmt.Errorf("Expected method %+v in %+v to have no parameters",
				method, targetInterface)
		}

		hasError := false
		if signature.Results().Len() == 2 {
			errType, ok := signature.Results().At(1).Type().(*types.Named)
			if !ok {
				return "", nil, fmt.Errorf("Expected %+v in %+v  to return an error",
					method, targetInterface)
			}

			if errType.Obj().Pkg() != nil {
				return "", nil, fmt.Errorf("Expected %+v in %+v  to return an error",
					method, targetInterface)
			}

			if errType.Obj().Name() != "error" {
				return "", nil, fmt.Errorf("Expected %+v in %+v  to return an error",
					method, targetInterface)
			}

			hasError = true
		}

		// Expect either one result or two results, the second one being an error
		if !(signature.Results().Len() == 1 || (signature.Results().Len() == 2 && hasError)) {
			return "", nil, fmt.Errorf("Expected method %+v in %+v to have one result and optional error",
				method, targetInterface)
		}

		isPointer := false
		realType := signature.Results().At(0).Type()
		var namedType *types.Named
		switch targetType := realType.(type) {
		case *types.Named:
			namedType = targetType
		case *types.Pointer:
			isPointer = true
			namedType = targetType.Elem().(*types.Named)
		default:
			return "", nil, fmt.Errorf("Type %+v is not a valid target", targetType)
		}

		targets = append(targets, &InjectionTarget{
			MethodName: method.Name(),
			Type:       realType,
			Name:       namedType,
			IsPointer:  isPointer,
			HasError:   hasError,
		})
	}

	return targetNamedType.Obj().Name(), targets, nil
}

func getNodesFromInterface(
	interfaceType *types.Interface,
	parent *resolutionNode,
) ([]*resolutionNode, error) {
	modulesMethod := typeutil.GetInterfaceMethod(interfaceType, modulesFunc)
	if modulesMethod == nil {
		return nil, nil
	}

	modulesMethodSignature := modulesMethod.Type().(*types.Signature)
	if modulesMethodSignature.Params().Len() > 0 {
		return nil, fmt.Errorf("Modules method %+v has arguments. Expected exactly 0", modulesMethod)
	}

	var nodes []*resolutionNode
	for i := 0; i < modulesMethodSignature.Results().Len(); i++ {
		nodes = append(nodes, &resolutionNode{
			parent:   parent,
			nodeType: modulesMethodSignature.Results().At(i).Type(),
		})
	}

	return nodes, nil
}

func extractBindings(
	node *structs.Interface,
) (map[string]*types.Named, error) {
	bindings := make(map[string]*types.Named)
	for i := 0; i < node.Type.NumMethods(); i++ {
		method := node.Type.Method(i)
		if !method.Exported() {
			continue
		}

		if _, ok := reservedMethods[method.Name()]; ok {
			continue
		}

		signature := method.Type().(*types.Signature)
		if signature.Params().Len() != 1 && signature.Results().Len() != 1 {
			return nil, fmt.Errorf("Expected method %+v in %+v to have one input and one output",
				method, node.Type)
		}

		interfaceName, ok := signature.Results().At(0).Type().(*types.Named)
		if !ok {
			return nil, fmt.Errorf("%+v was not named in %+v", signature.Params().At(0).Type(), node)
		}

		interfaceID := typeutil.IDFromNamed(interfaceName)
		if _, ok := bindings[interfaceID]; ok {
			return nil, fmt.Errorf("Found duplicate binding for %+v in %+v", interfaceName, node)
		}

		var implementationName *types.Named
		implementationType := signature.Params().At(0).Type()
		switch actualType := implementationType.(type) {
		case *types.Pointer:
			name, ok := actualType.Elem().(*types.Named)
			if !ok {
				return nil, fmt.Errorf("Expecting %+v to be a struct in %+v", implementationName, node)
			}
			implementationName = name
		case *types.Named:
			implementationName = actualType
		default:
			return nil, fmt.Errorf("%+v is not a pointer or a named type", implementationType)
		}

		bindings[interfaceID] = implementationName
	}

	return bindings, nil
}

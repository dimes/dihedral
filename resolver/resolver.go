// Package resolver handles resolving dependencies
package resolver

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/dimes/di/structs"
	"github.com/dimes/di/typeutil"
	"github.com/pkg/errors"
)

const (
	modulesFunc = "Modules"
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
	Type       *types.Named
	IsPointer  bool
}

// ResolvedType represents a type that has been resolved via a module.
// There are two possible ways a type can be resolved:
// 1. Being bound via an interface module
// 2. Being provided via a struct module
type ResolvedType struct {
	Module    *structs.Struct
	Method    *types.Func
	Name      *types.Named
	IsPointer bool
}

// ResolveComponentModules resolves the modules for the component interface.
// The return types are:
// - List of struct modules (used to provide concrete types)
// - List of interface modules (used to bind interfaces to implementations)
func ResolveComponentModules(
	fileSet *token.FileSet,
	componentInterface *structs.Interface,
) (
	[]*InjectionTarget,
	map[string]*ResolvedType,
	map[string]*structs.Struct,
	error,
) {
	targets, err := getTargetsFromInterface(componentInterface.Type)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "Error getting targets for %+v", componentInterface)
	}

	stack, err := getNodesFromInterface(componentInterface.Type, nil)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "Error getting modules for %+v", componentInterface)
	}

	seen := make(map[string]struct{})
	providers := make(map[string]*ResolvedType)
	bindings := make(map[string]*structs.Struct)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch typedNode := node.nodeType.(type) {
		case *types.Named:
			nodeInterface, ok := typedNode.Underlying().(*types.Interface)
			if !ok {
				return nil, nil, nil, fmt.Errorf("Expected node %+v to be pointer or interface", typedNode)
			}

			id := typeutil.IDFromNamed(typedNode)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}

			nodeModules, err := getNodesFromInterface(nodeInterface, node)
			if err != nil {
				return nil, nil, nil, errors.Wrapf(err, "Error getting dependencies for %+v", nodeInterface)
			}

			stack = append(stack, nodeModules...)

			bindingInterface := &structs.Interface{
				Name: typedNode,
				Type: nodeInterface,
			}

			moduleBindings, err := extractBindings(bindingInterface)
			if err != nil {
				return nil, nil, nil, errors.Wrapf(err, "Error extracting bindings in %+v", nodeInterface)
			}

			for id, boundStruct := range moduleBindings {
				if _, ok := bindings[id]; ok {
					return nil, nil, nil, fmt.Errorf("Binding %+v seen twice", id)
				}

				if _, ok := providers[id]; ok {
					return nil, nil, nil, fmt.Errorf("Binding %+v seen twice", id)
				}

				bindings[id] = boundStruct
			}
		case *types.Pointer:
			namedNode, ok := typedNode.Elem().(*types.Named)
			if !ok {
				return nil, nil, nil, fmt.Errorf("Expected pointer %+v to point to named element",
					typedNode)
			}

			id := typeutil.IDFromNamed(namedNode)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}

			structNode, ok := namedNode.Underlying().(*types.Struct)
			if !ok {
				return nil, nil, nil, fmt.Errorf("Expected pointer %+v to point to a struct", typedNode)
			}

			module := &structs.Struct{
				Name: namedNode,
				Type: structNode,
			}

			packageName := namedNode.Obj().Pkg().Name()
			packages, err := parser.ParseDir(fileSet, packageName, nil, 0)
			if err != nil {
				return nil, nil, nil, errors.Wrapf(err, "Error parsing package %s", packageName)
			}

			for _, astPkg := range packages {
				var files []*ast.File
				for _, file := range astPkg.Files {
					files = append(files, file)
				}

				info := &types.Info{
					Defs: make(map[*ast.Ident]types.Object),
				}

				conf := types.Config{
					Importer: importer.ForCompiler(fileSet, "source", nil),
				}

				_, err := conf.Check(packageName, fileSet, files, info)
				if err != nil {
					return nil, nil, nil, errors.Wrapf(err, "Error getting definitions for package %s",
						packageName)
				}

				for identifier, definition := range info.Defs {
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

					if signature.Results().Len() != 1 {
						return nil, nil, nil, fmt.Errorf("Expecting exactly one result from %+v", signature)
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
						return nil, nil, nil, fmt.Errorf("Result %+v is an unsupported type", result)
					}

					resultID := typeutil.IDFromNamed(resultName)
					if _, ok := bindings[resultID]; ok {
						return nil, nil, nil, fmt.Errorf("Binding %+v seen twice", resultID)
					}

					if _, ok := providers[resultID]; ok {
						return nil, nil, nil, fmt.Errorf("Binding %+v seen twice", resultID)
					}

					resolvedType := &ResolvedType{
						Module:    module,
						Method:    funcDefinition,
						Name:      resultName,
						IsPointer: isPointer,
					}

					providers[resultID] = resolvedType
				}
			}
		default:
			return nil, nil, nil, fmt.Errorf("%+v is not a recognized module type", typedNode)
		}
	}

	return targets, providers, bindings, nil
}

func getTargetsFromInterface(
	interfaceType *types.Interface,
) ([]*InjectionTarget, error) {
	targets := make([]*InjectionTarget, 0)
	for i := 0; i < interfaceType.NumMethods(); i++ {
		method := interfaceType.Method(i)
		if !method.Exported() {
			continue
		}

		if _, ok := reservedMethods[method.Name()]; ok {
			continue
		}

		signature := method.Type().(*types.Signature)
		if signature.Params().Len() > 0 || signature.Results().Len() != 1 {
			return nil, fmt.Errorf("Expected method %+v in %+v to have no arguments and one result",
				method, interfaceType)
		}

		isPointer := false
		var namedType *types.Named
		switch targetType := signature.Results().At(0).Type().(type) {
		case *types.Named:
			namedType = targetType
		case *types.Pointer:
			isPointer = true
			namedType = targetType.Elem().(*types.Named)
		default:
			return nil, fmt.Errorf("Type %+v is not a valid target", targetType)
		}

		targets = append(targets, &InjectionTarget{
			MethodName: method.Name(),
			Type:       namedType,
			IsPointer:  isPointer,
		})
	}

	return targets, nil
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
) (map[string]*structs.Struct, error) {
	bindings := make(map[string]*structs.Struct)
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

		implementationName, ok := signature.Params().At(0).Type().(*types.Named)
		if !ok {
			return nil, fmt.Errorf("%+v was not named in %+v", signature.Params().At(0).Type(), node)
		}

		implementationType, ok := implementationName.Underlying().(*types.Struct)
		if !ok {
			return nil, fmt.Errorf("Expecting %+v to be a struct in %+v", implementationName, node)
		}

		bindings[interfaceID] = &structs.Struct{
			Name: implementationName,
			Type: implementationType,
		}
	}

	return bindings, nil
}

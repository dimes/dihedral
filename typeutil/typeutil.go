// Package typeutil contains helepr methods for interaction with ast types
package typeutil

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/dimes/di/structs"
	"github.com/pkg/errors"
)

// IDFromNamed returns a unique string for the given name
func IDFromNamed(name *types.Named) string {
	return name.Obj().Pkg().Path() + "." + name.Obj().Name()
}

// FindInterface finds the given interface name in the given package. Returns
// nil if no interface is found.
func FindInterface(
	fileSet *token.FileSet,
	packageName string,
	interfaceName string,
) (*structs.Interface, error) {
	packages, err := parser.ParseDir(fileSet, packageName, nil, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Error parsing package %s", packageName)
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
			return nil, errors.Wrapf(err, "Error getting defs for package %s", packageName)
		}

		for identifier, definition := range info.Defs {
			if !identifier.IsExported() {
				continue
			}

			if identifier.Name != interfaceName {
				continue
			}

			namedType, ok := definition.Type().(*types.Named)
			if !ok {
				return nil, fmt.Errorf("Type %+v is not a named type", definition)
			}

			interfaceType, ok := namedType.Underlying().(*types.Interface)
			if !ok {
				return nil, fmt.Errorf("%s in %s is not an interface", interfaceName, packageName)
			}

			return &structs.Interface{
				Name: namedType,
				Type: interfaceType,
			}, nil
		}
	}

	return nil, nil
}

// GetInterfaceMethod returns the method with given name in the interface, or nil
func GetInterfaceMethod(interfaceType *types.Interface, methodName string) *types.Func {
	for i := 0; i < interfaceType.NumMethods(); i++ {
		method := interfaceType.Method(i)
		if !method.Exported() {
			continue
		}

		if method.Name() == methodName {
			return method
		}
	}

	return nil
}

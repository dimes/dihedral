// Package typeutil contains helepr methods for interaction with ast types
package typeutil

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"

	"github.com/dimes/dihedral/structs"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
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
	config := &packages.Config{
		Mode: packages.LoadSyntax,
	}
	pkgs, err := packages.Load(config, packageName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading package %s", packageName)
	}

	for _, astPkg := range pkgs {
		for identifier, definition := range astPkg.TypesInfo.Defs {
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

// HasFieldOfType returns true if the given struct has a non-exported field of type
// fieldType
func HasFieldOfType(
	targetStruct *types.Struct,
	fieldType reflect.Type,
) bool {
	for i := 0; i < targetStruct.NumFields(); i++ {
		field := targetStruct.Field(i)
		if field.Exported() {
			continue
		}

		namedType, ok := field.Type().(*types.Named)
		if !ok {
			continue
		}

		if fieldType.PkgPath() != namedType.Obj().Pkg().Path() {
			continue
		}

		if fieldType.Name() != namedType.Obj().Name() {
			continue
		}

		return true
	}

	return false
}

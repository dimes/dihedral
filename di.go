package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/token"

	"github.com/dimes/di/gen"
	"github.com/dimes/di/resolver"
	"github.com/dimes/di/typeutil"
)

const (
	modulesFunc = "Modules"
)

func main() {
	var packageName string
	var componentName string

	flag.StringVar(&packageName, "package", "test", "The name of the package containing the component")
	flag.StringVar(&componentName, "component", "SomeType", "The name of the component")
	flag.Parse()

	fileSet := token.NewFileSet()

	componentInterface, err := typeutil.FindInterface(fileSet, packageName, componentName)
	if err != nil {
		panic(err)
	}

	if componentInterface == nil {
		panic("Component interface not found")
	}

	targets, providers, bindings, err := resolver.ResolveComponentModules(fileSet, componentInterface)
	if err != nil {
		panic(err)
	}

	component, err := gen.NewGeneratedComponent(componentName, targets, providers, bindings)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found targets: %+v\n", targets)
	fmt.Printf("Found providers: %+v\n", providers)
	fmt.Printf("Found bindings: %+v\n", bindings)
	fmt.Printf("Generated component: %+v\n", component)

	generatedSource := component.ToSource("di")
	for name, file := range generatedSource {
		fmt.Println("------")
		fmt.Println(name + ".go")
		fmt.Println(file)
	}
}

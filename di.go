package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

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
	var outputDir string

	flag.StringVar(&packageName, "package", "", "The name of the package containing the component")
	flag.StringVar(&componentName, "component", "MyComponent", "The name of the component")
	flag.StringVar(&outputDir, "output", "di", "The directory to output generated source to")
	flag.Parse()

	if componentName == "" {
		panic("-component must be set")
	}

	if packageName == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		cmd := exec.Command("go", "list")
		cmd.Env = append(append([]string{}, os.Environ()...), "PWD="+workingDir)
		cmd.Dir = workingDir
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}

		packageName = strings.TrimSpace(stdout.String())
	}

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

	fmt.Printf("Found targets: %+v\n", targets)
	fmt.Printf("Found providers: %+v\n", providers)
	fmt.Printf("Found bindings: %+v\n", bindings)

	component, err := gen.NewGeneratedComponent(componentName, targets, providers, bindings)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated component: %+v\n", component)

	generatedSource := component.ToSource("di")
	os.MkdirAll(outputDir, os.ModePerm)
	for name, file := range generatedSource {
		if err := ioutil.WriteFile(
			path.Join(outputDir, name+".go"),
			[]byte(file),
			os.ModePerm,
		); err != nil {
			panic(err)
		}
	}
}

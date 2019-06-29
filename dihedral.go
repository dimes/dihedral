package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dimes/dihedral/gen"
	"github.com/dimes/dihedral/resolver"
	"github.com/dimes/dihedral/typeutil"
)

const (
	modulesFunc = "Modules"
)

func main() {
	var packageName string
	var definitionName string
	var outputDir string

	flag.StringVar(&packageName, "package", "", "The name of the package containing the component")
	flag.StringVar(&definitionName, "definition", "", "The name of the definition interface")
	flag.StringVar(&outputDir, "output", "digen", "The directory to output generated source to")
	flag.Parse()

	if definitionName == "" {
		panic("-definition must be set")
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

	definitionInterface, err := typeutil.FindInterface(fileSet, packageName, definitionName)
	if err != nil {
		panic(err)
	}

	if definitionInterface == nil {
		panic("Definition interface not found")
	}

	result, err := resolver.ResolveComponentModules(fileSet, definitionInterface)
	if err != nil {
		panic(err)
	}

	targetInterfaceName := result.TargetInterfaceName
	targets := result.Targets
	providers := result.Providers
	bindings := result.Bindings

	fmt.Printf("Found target interface %s\n", targetInterfaceName)
	fmt.Printf("Found targets: %+v\n", targets)
	fmt.Printf("Found providers: %+v\n", providers)
	fmt.Printf("Found bindings: %+v\n", bindings)

	component, err := gen.NewGeneratedComponent(targetInterfaceName, targets, providers, bindings)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated component: %+v\n", component)

	generatedSource := component.ToSource(filepath.Base(outputDir))
	os.MkdirAll(outputDir, os.ModePerm)
	for name, file := range generatedSource {
		formatted, err := format.Source([]byte(file))
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(
			path.Join(outputDir, name+".go"),
			formatted,
			os.ModePerm,
		); err != nil {
			panic(err)
		}
	}
}

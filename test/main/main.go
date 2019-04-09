package main

import (
	"fmt"

	"github.com/dimes/di/test"

	"github.com/dimes/di/test/di"
)

func main() {
	component := di.NewMyComponent(&test.Module{Greeting: "Salutations"})
	fmt.Printf("%+v\n", component.Target())
	fmt.Printf("%+v\n", component.Target().Greeting)
	fmt.Printf("%+v\n", component.Target().MyInterface.TestMethod())
}

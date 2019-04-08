package main

import (
	"fmt"

	"github.com/dimes/di/test/di"
)

func main() {
	component := di.NewMyComponent()
	fmt.Printf("%+v\n", component.Target())
	fmt.Printf("%+v\n", component.Target().MyInterface.TestMethod())
}

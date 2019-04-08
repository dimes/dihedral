package main

import (
	"fmt"

	"github.com/dimes/di/di"
)

func main() {
	component := di.NewMyComponent()
	fmt.Printf("%+v\n", component.Target())
}

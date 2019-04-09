package main

import (
	"fmt"

	"github.com/dimes/di/example/bindings/di"
)

func main() {
	component := di.NewServiceComponent()
	service := component.GetService()

	if err := service.SetValueInDBStore("Hello World!"); err != nil {
		panic(err)
	}

	fmt.Println(service.GetValueFromDBStore())
}

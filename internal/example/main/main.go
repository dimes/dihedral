package main

import (
	"fmt"

	"github.com/dimes/dihedral/internal/example/bindings/digen"
	"github.com/dimes/dihedral/internal/example/dbstore"
)

func main() {
	component := digen.NewServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})

	timeout := component.GetServiceTimeout()
	fmt.Println("Service timeout is", timeout)

	service, err := component.GetService()
	if err != nil {
		panic(err)
	}

	if err := service.SetValueInDBStore("World!"); err != nil {
		panic(err)
	}

	fmt.Println(service.GetValueFromDBStore())
}

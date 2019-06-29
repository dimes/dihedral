package main

import (
	"fmt"

	"github.com/dimes/dihedral/internal/example/bindings"
	"github.com/dimes/dihedral/internal/example/bindings/digen"
	"github.com/dimes/dihedral/internal/example/dbstore"
)

func main() {
	var component bindings.ServiceComponent
	component = digen.NewDihedralServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})

	timeout, err := component.GetServiceTimeout()
	if err != nil {
		panic(err)
	}

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

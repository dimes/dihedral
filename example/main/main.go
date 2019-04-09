package main

import (
	"fmt"

	"github.com/dimes/di/example/bindings/di"
	"github.com/dimes/di/example/dbstore"
)

func main() {
	component := di.NewServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})
	service := component.GetService()

	if err := service.SetValueInDBStore("World!"); err != nil {
		panic(err)
	}

	fmt.Println(service.GetValueFromDBStore())
}

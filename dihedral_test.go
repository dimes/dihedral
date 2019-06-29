package main

import (
	"testing"

	"github.com/dimes/dihedral/internal/example/bindings/digen"
	"github.com/dimes/dihedral/internal/example/dbstore"
	"github.com/stretchr/testify/assert"
)

func TestExampleInjection(t *testing.T) {
	component := digen.NewDihedralServiceComponent(&dbstore.DBProviderModule{
		Prefix: "Hello",
	})

	service, err := component.GetService()
	assert.NoError(t, err)
	assert.NoError(t, service.SetValueInDBStore("World!"))
	assert.Equal(t, "Hello World!", service.GetValueFromDBStore())

	serviceTimeout, err := component.GetServiceTimeout()
	assert.NoError(t, err)
	assert.Equal(t, 5000000000, serviceTimeout)
}

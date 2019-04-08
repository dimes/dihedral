package di

import (
	di_import_1 "github.com/dimes/di/test"
)

type GeneratedComponent struct {
	test_Module *di_import_1.Module
}

func NewMyComponent() *GeneratedComponent {
	return &GeneratedComponent{
		test_Module: &di_import_1.Module{},
	}
}
func (generatedComponent *GeneratedComponent) Target() *di_import_1.MyTarget {
	return factory_test_MyTarget(generatedComponent)
}

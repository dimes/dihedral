package di

import target_pkg "github.com/dimes/di/test"

func factory_test_MyImplementation(generatedComponent *GeneratedComponent) *target_pkg.MyImplementation {
	target := &target_pkg.MyImplementation{}
	return target
}

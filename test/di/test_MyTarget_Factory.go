package di

import target_pkg "github.com/dimes/di/test"

func factory_test_MyTarget(generatedComponent *GeneratedComponent) *target_pkg.MyTarget {
	target := &target_pkg.MyTarget{}
	target.Greeting = generatedComponent.provides_test_Greeting()
	target.Timeout = generatedComponent.provides_test_ServiceTimeout()
	target.MyInterface = factory_test_MyImplementation(generatedComponent)
	return target
}

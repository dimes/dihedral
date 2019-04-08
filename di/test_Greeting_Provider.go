package di

import target_pkg "github.com/dimes/di/test"

func (generatedComponent *GeneratedComponent) provides_test_Greeting() target_pkg.Greeting {
	return generatedComponent.test_Module.ProvidesGreeting()
}

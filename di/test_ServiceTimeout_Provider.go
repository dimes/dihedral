package di

import target_pkg "github.com/dimes/di/test"

func (generatedComponent *GeneratedComponent) provides_test_ServiceTimeout() target_pkg.ServiceTimeout {
	return generatedComponent.test_Module.ProvideServiceTimeout()
}

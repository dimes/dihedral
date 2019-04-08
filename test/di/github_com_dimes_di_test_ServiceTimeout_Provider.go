package di
import target_pkg "github.com/dimes/di/test"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_di_test_ServiceTimeout() target_pkg.ServiceTimeout {
	return generatedComponent.github_com_dimes_di_test_Module.ProvideServiceTimeout(
	)
}

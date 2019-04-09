package di
import target_pkg "github.com/dimes/di/example"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_di_example_ServiceTimeout() target_pkg.ServiceTimeout {
	return generatedComponent.github_com_dimes_di_example_bindings_ServiceModule.ProvidesServiceTimeout(
	)
}

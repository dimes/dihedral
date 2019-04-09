package di
import target_pkg "github.com/dimes/di/example"
func factory_github_com_dimes_di_example_Service(generatedComponent *GeneratedComponent) *target_pkg.Service {
	target := &target_pkg.Service{}
	target.ServiceTimeout = generatedComponent.provides_github_com_dimes_di_example_ServiceTimeout()
	target.DBStore = factory_github_com_dimes_di_example_dbstore_TableDBStore(generatedComponent)
	return target
}

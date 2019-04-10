package digen
import target_pkg "github.com/dimes/dihedral/example"
func factory_github_com_dimes_dihedral_example_Service(generatedComponent *GeneratedComponent) *target_pkg.Service {
	target := &target_pkg.Service{}
	target.ServiceTimeout = generatedComponent.provides_github_com_dimes_dihedral_example_ServiceTimeout()
	target.DBStore = factory_github_com_dimes_dihedral_example_dbstore_MemoryDBStore(generatedComponent)
	return target
}

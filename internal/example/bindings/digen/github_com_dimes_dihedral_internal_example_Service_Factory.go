package digen
import target_pkg "github.com/dimes/dihedral/internal/example"
func factory_github_com_dimes_dihedral_internal_example_Service(generatedComponent *GeneratedComponent) *target_pkg.Service {
	target := &target_pkg.Service{}
	target.DBStore = factory_github_com_dimes_dihedral_internal_example_dbstore_MemoryDBStore(generatedComponent)
	target.ServiceTimeout = generatedComponent.provides_github_com_dimes_dihedral_internal_example_ServiceTimeout()
	return target
}

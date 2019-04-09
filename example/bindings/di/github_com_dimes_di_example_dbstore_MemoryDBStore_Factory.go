package di
import target_pkg "github.com/dimes/di/example/dbstore"
func factory_github_com_dimes_di_example_dbstore_MemoryDBStore(generatedComponent *GeneratedComponent) *target_pkg.MemoryDBStore {
	target := &target_pkg.MemoryDBStore{}
	target.Prefix = generatedComponent.provides_github_com_dimes_di_example_dbstore_Prefix()
	return target
}

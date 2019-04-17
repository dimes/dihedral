package digen
import target_pkg "github.com/dimes/dihedral/internal/example/dbstore"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_dihedral_internal_example_dbstore_Prefix() target_pkg.Prefix {
	return generatedComponent.github_com_dimes_dihedral_internal_example_dbstore_DBProviderModule.ProvidesPrefix(
	)
}

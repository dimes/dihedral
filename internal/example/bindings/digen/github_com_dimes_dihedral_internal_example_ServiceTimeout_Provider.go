package digen
import target_pkg "github.com/dimes/dihedral/internal/example"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_dihedral_internal_example_ServiceTimeout() target_pkg.ServiceTimeout {
	return generatedComponent.github_com_dimes_dihedral_internal_example_bindings_ServiceModule.ProvidesServiceTimeout(
	)
}

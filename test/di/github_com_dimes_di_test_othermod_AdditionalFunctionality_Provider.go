package di
import target_pkg "github.com/dimes/di/test/othermod"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_di_test_othermod_AdditionalFunctionality() *target_pkg.AdditionalFunctionality {
	return generatedComponent.github_com_dimes_di_test_othermod_Module.ProvidesAdditionalFunctionality(
	)
}

package di
import target_pkg "github.com/dimes/di/test"
func (generatedComponent *GeneratedComponent) provides_github_com_dimes_di_test_Greeting() target_pkg.Greeting {
	return generatedComponent.github_com_dimes_di_test_Module.ProvidesGreeting(
	)
}

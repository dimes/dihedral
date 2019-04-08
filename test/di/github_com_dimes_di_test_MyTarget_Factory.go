package di
import target_pkg "github.com/dimes/di/test"
func factory_github_com_dimes_di_test_MyTarget(generatedComponent *GeneratedComponent) *target_pkg.MyTarget {
	target := &target_pkg.MyTarget{}
	target.Greeting = generatedComponent.provides_github_com_dimes_di_test_Greeting()
	target.Timeout = generatedComponent.provides_github_com_dimes_di_test_ServiceTimeout()
	target.MyInterface = factory_github_com_dimes_di_test_MyImplementation(generatedComponent)
	target.Additional = generatedComponent.provides_github_com_dimes_di_test_othermod_AdditionalFunctionality()
	return target
}

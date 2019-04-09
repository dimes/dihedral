package di
import (
	di_import_2 "github.com/dimes/di/test"
	di_import_1 "github.com/dimes/di/test/othermod"
)
type GeneratedComponent struct {
	github_com_dimes_di_test_othermod_Module *di_import_1.Module
	github_com_dimes_di_test_Module *di_import_2.Module
}
func NewMyComponent(
	github_com_dimes_di_test_Module *di_import_2.Module,
) *GeneratedComponent {
	 return &GeneratedComponent{
		github_com_dimes_di_test_othermod_Module: &di_import_1.Module{},
		github_com_dimes_di_test_Module: github_com_dimes_di_test_Module,
	}
}
func (generatedComponent *GeneratedComponent) Target() *di_import_2.MyTarget {
	return factory_github_com_dimes_di_test_MyTarget(generatedComponent)
}

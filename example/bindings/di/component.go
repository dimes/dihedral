package di
import (
	di_import_1 "github.com/dimes/di/example/bindings"
	di_import_2 "github.com/dimes/di/example"
)
type GeneratedComponent struct {
	github_com_dimes_di_example_bindings_ServiceModule *di_import_1.ServiceModule
}
func NewServiceComponent(
) *GeneratedComponent {
	 return &GeneratedComponent{
		github_com_dimes_di_example_bindings_ServiceModule: &di_import_1.ServiceModule{},
	}
}
func (generatedComponent *GeneratedComponent) GetService() *di_import_2.Service {
	return factory_github_com_dimes_di_example_Service(generatedComponent)
}

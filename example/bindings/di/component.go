package di
import (
	di_import_1 "github.com/dimes/di/example/dbstore"
	di_import_2 "github.com/dimes/di/example/bindings"
	di_import_3 "github.com/dimes/di/example"
)
type GeneratedComponent struct {
	github_com_dimes_di_example_dbstore_DBProviderModule *di_import_1.DBProviderModule
	github_com_dimes_di_example_bindings_ServiceModule *di_import_2.ServiceModule
}
func NewServiceComponent(
	github_com_dimes_di_example_dbstore_DBProviderModule *di_import_1.DBProviderModule,
) *GeneratedComponent {
	 return &GeneratedComponent{
		github_com_dimes_di_example_dbstore_DBProviderModule: github_com_dimes_di_example_dbstore_DBProviderModule,
		github_com_dimes_di_example_bindings_ServiceModule: &di_import_2.ServiceModule{},
	}
}
func (generatedComponent *GeneratedComponent) GetService() *di_import_3.Service {
	return factory_github_com_dimes_di_example_Service(generatedComponent)
}

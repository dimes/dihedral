# dihedral

**dihedral** is a compile-time injection framework for Go.

# Getting started

    > go get github.com/dimes/dihedral

Create a type you want injected

    type ServiceEndpoint string  // Name this string "ServiceEndpoint"
    type Service struct {
        inject  embeds.Inject    // Auto-inject this struct 
        Endpoint ServiceEndpoint // Inject a string with name "ServiceEndpoint"
    }

Create a module to provide non-injected dependencies

    // Each public method on this struct provides a type
    type ServiceModule struct {}
    func (s *ServiceModule) ProvidesServiceEndpoint() ServiceEndpoint {
        return ServiceEndpoint("http://hello.world")
    }

Create a component as the root of the dependency injection

    interface ServiceComponent {
        Modules() *MyModule      // Tells dihedral which modules to include
        InjectService() *Service // Tells dihedral the root of the DI graph
    }

Generate the bindings

    > dihedral -component ServiceComponent

Use the bindings

    func main() {
        // dihedral generates the digen package
        component := digen.ServiceComponent()
        service := component.InjectService()
        fmt.Println(string(injected.Endpoint)) # Prints "http://hello.world"
    }

See the [example](example/) for a more detailed overview.

## Concepts 

**dihedral** is powered by four simple concepts.

#### 1. Provider Modules

Provider modules are structs with associated "provider" methods, which is an exported method that returns an instance of some type. These methods are called directly at runtime when the injection is performed. Provider functions can take parameters that are injected.

```
type MyProviderModule struct {}
func (m *MyProviderModule) ProvidesSQLDatabase(tableName TableName) *SQLDatabase {
    return &SQLDatabase{tableName: tableName}
}
```

#### 2. Binding Modules

Binding modules are interfaces that "bind" implementations to an interface. Each method on the interface is treated as a type binding and is expected to have one parameter and one return type. The return type is an interface, and the parameter type is a struct that implements that interface. During injection, instances of that interface will be provided by the bound implementation. 

Binding modules also have a special `Modules()` method. This method takes no parameters and returns a list of other modules to use. This is useful for encapsulating logic in modules that can then be used as an entire unit.

```
type MyBindingModule interface {
    Modules() (*MyProviderModule, MyOtherBindingModule)
    BindsDatabase(impl *SQLDatabase) Database
}
```

#### 3. Injectable structs

**dihedral** can automatically inject structs. This reduces a huge amount of boilerplate and keeps provider modules very small. To inject a struct, it must contain a non-exported field of type `embeds.Inject`. By default, all exported fields are injected. To skip a field, add the tag `di:"-"`.

```
type Service struct {
    inject       embeds.Inject
    ServiceDB    Database
    RequestCount int           `di:"-"`
}
```

#### 4. Components

Components tie everything together. Ultimately, **dihedral** generates a component struct and constructor function. The methods on a component return fully injected types. Like binding modules, component interfaces also include a `Modules()` method that specifies which modules to use for injection.

```
type ServiceComponent interface {
    Modules() (*MyProviderModule, MyBindingModule)
    InjectService() *Service
}
```

## Code Generation

Code can be generated using the CLI:

    > dihedral -component ServiceComponent

This can also be used with `go generate` by putting the following at the top of the file containing the component:

    //go:generate dihedral -component ServiceComponent

The generated code is located in the `digen` folder. If the component is named `ServiceComponent`, then **dihedral** generates a function named `func NewServiceComponent() *GeneratedServiceComponent` that contains the requested injections. 

```
func main() {
    component := digen.NewServiceComponent()
    service := component.InjectService()
}
```

## Runtime Values

In the previous examples, all of the required injection information was known at compile time. Often, certain injection parameters are only known at runtime, e.g. configuration values such as the endpoint to use, the database table name, etc. **dihedral** solves this problem by allowing modules to be provided at runtime.

For instance, let's say we have the following module that provides configuration:

```
type Config struct {
    TableName string
}

type ConfigModule struct {
    provided embeds.ProvidedModule // Tells dihedral to get this at runtime
    Config *Config
}

type TableName string
func (c *ConfigModule) ProvidesTableName() TableName {
    return TableName(c.Config.TableName)
}
```

By adding a non-exported field of type `embeds.ProvidedModule` to the module, **dihedral** will add the module as a parameter to the generated component. Instead of

    func NewServiceComponent() *GeneratedComponent

The generated function will look like 

    func NewServiceComponent(module *ConfigModule) *GeneratedComponent

It can be used like this

    func main() {
        config := &Config{ TableName: "test-table" }
        module := &ConfigModule { Config: config }
        component := digen.NewServiceComponent(module)
        service := component.InjectService()
    }

### Differences from Wire

Wire, Google's injection framework, is another compile time framework for Go. Both frameworks are inspired
by Dagger. **dihedral** differs from Wire in that **dihedral** focuses on auto-injected components and self-contained modules, whereas Wire focuses more on type registration via provider functions. **dihedral** also leverages struct receivers for better organization of runtime provided types. These features make **dihedral** a pleasure to work with. 

**dihedral**'s component structure also enables one to have multiple injected components that share modules. The type annotation system allows for auto-injected components, provided modules, and, in the future, sub-components and have a different scope than the parent component.

---
layout: sidebar
title: Dihedral
---

## Provider Modules

Provider modules are structs with associated "provider" methods, which is an exported method that returns an instance of some type. These methods are called directly at runtime when the injection is performed. Provider functions can take parameters that are injected.

```
type MyProviderModule struct {}
func (m *MyProviderModule) ProvidesSQLDatabase(tableName TableName) *SQLDatabase {
    return &SQLDatabase{tableName: tableName}
}
```

### Runtime Values

Runtime values can be provided by constructing module instances at runtime and using them as constructor parameters in the generated component factory function.

For instance, consider the following struct that provides configuration:

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

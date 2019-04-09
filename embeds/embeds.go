package embeds

// Inject is an empty struct that can be added as a non-exported parameter
// in other structs to indicate that injection should happen automatically
type Inject struct {
}

// ProvidedModule is an empty struct that can be added as a non-exported
// parameter to a module to indicate it should be a parameter of the component
type ProvidedModule struct {
}

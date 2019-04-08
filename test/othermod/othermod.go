package othermod

type AdditionalFunctionality struct {
}

type Module struct {
}

func (m *Module) ProvidesAdditionalFunctionality() *AdditionalFunctionality {
	return &AdditionalFunctionality{}
}

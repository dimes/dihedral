package test

import (
	"io"

	"github.com/dimes/di/embeds"
	"github.com/dimes/di/something"
)

func A() {

}

type Module struct {
}

func (m *Module) SomeBinding() MyString {
	return MyString("Hello")
}

func (m *Module) OtherBinding(input MyString) BindingModule {
	return nil
}

type BindingModule interface {
	BindsR(impl R) io.Reader
}

type SomeType interface {
	Modules() (*Module, *something.Module, BindingModule)
	Target() *MyTarget
}

type R struct {
}

func (r *R) Read(p []byte) (n int, err error) {
	return 0, nil
}

type MyString string
type MyTarget struct {
	inject embeds.Inject

	A      MyString `di:"tag"`
	TheMod BindingModule
}

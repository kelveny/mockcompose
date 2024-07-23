package foo

import "fmt"

type Foo interface {
	Foo() string
	Bar() bool
}

type foo struct {
	name string
}

var _ Foo = (*foo)(nil)
var _ Foo = (*dummyFoo)(nil)

func (f *foo) Foo() string {
	if f.Bar() {
		return "Overriden with Bar"
	}

	return f.name
}

func (f *foo) Bar() bool {
	return f.name == "bar"
}

// dummyInner/dummyFoo is used for callee detection test
//  1. a peer method callee
//  2. function callee from within the package
//  3. function callee from other package
//  4. callee from function pointer
//  5. callee from function pointer of embedded type
type dummyInner struct {
	fp func()
}
type dummyFoo struct {
	name string

	in   dummyInner
	fptr func()
}

func (f *dummyFoo) Foo() string {
	if f.Bar() { // peer method callee
		return "Overriden with Bar"
	}

	f.fptr()  // callee via field pointer
	f.in.fp() // callee via embedded inner filed pointer

	dummy() // function callee from within the package

	fmt.Printf("dummy") // function callee from other package

	return f.name
}

func (f *dummyFoo) Bar() bool {
	return f.name == "bar"
}

func dummy() {
}

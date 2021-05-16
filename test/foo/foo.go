package foo

type Foo interface {
	Foo() string
	Bar() bool
}

type foo struct {
	name string
}

var _ Foo = (*foo)(nil)

func (f *foo) Foo() string {
	if f.Bar() {
		return "Overriden with Bar"
	}

	return f.name
}

func (f *foo) Bar() bool {
	if f.name == "bar" {
		return true
	}

	return false
}

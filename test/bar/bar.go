package bar

import (
	"fmt"
	"math/rand"
	"time"
)

type fooBar struct {
	name string
}

//go:generate mockcompose -n fooBarMock -c fooBar -real FooBar,this -real BarFoo,this:.

func (f *fooBar) FooBar() string {
	if f.order()%2 == 0 {
		fmt.Printf("ordinal order\n")

		return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
	}

	fmt.Printf("reverse order\n")
	return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
}

func (f *fooBar) BarFoo() string {
	if order()%2 == 0 {
		fmt.Printf("ordinal order\n")

		return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
	}

	fmt.Printf("reverse order\n")
	return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
}

func (f *fooBar) Foo() string {
	return "Foo"
}

func (f *fooBar) Bar() string {
	return "Bar"
}

func (f *fooBar) order() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

func order() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

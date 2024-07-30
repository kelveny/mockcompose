//go:generate mockcompose -n testFoo -c foo -real Foo,this:.:fmt
//go:generate mockcompose -n FooMock -i Foo
package foo

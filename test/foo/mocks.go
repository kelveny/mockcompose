//go:generate mockcompose -n testFoo -c foo -real Foo -mock Bar
//go:generate mockcompose -n FooMock -i Foo
package foo

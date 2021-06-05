//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
//go:generate mockcompose -n clonedFuncs -real "functionThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock" -real "functionThatUsesGlobalFunction,fmt=fmtMock" -real "functionThatUsesMultileGlobalFunctions2,fmt=fmtMock"
package clonefn

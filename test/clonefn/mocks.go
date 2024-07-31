//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
//go:generate mockcompose -n clonedFuncs -real "functionThatUsesMultileGlobalFunctions,fmt:json" -real "functionThatUsesGlobalFunction,fmt" -real "functionThatUsesMultileGlobalFunctions2,fmt"
package clonefn

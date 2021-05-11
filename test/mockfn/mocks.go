//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
//go:generate mockcompose -n mockSampleClz -c sampleClz -real "methodThatUsesGlobalFunction,fmt=fmtMock"
//go:generate mockcompose -n mockSampleClz2 -c sampleClz -real "methodThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock"
//go:generate mockcompose -n mockSampleClz3 -c sampleClz -real "methodThatUsesMultileGlobalFunctions,fmt=fmtMock"
package mockfn

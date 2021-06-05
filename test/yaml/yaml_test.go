package yaml

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var jsonMock *mockJson = &mockJson{}
var fmtMock *mockFmt = &mockFmt{}

func TestSampleClz(t *testing.T) {
	assert := require.New(t)

	// setup function mocks
	jsonMock.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	fmtMock.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	// inside mockSampleClz.methodThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	sc := mockSampleClz{}
	assert.True(sc.methodThatUsesGlobalFunction("format", "value") == "mocked Sprintf")

	// inside mockSampleClz2.methodThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	sc2 := mockSampleClz2{}
	assert.True(sc2.methodThatUsesMultileGlobalFunctions("format", "value") == "mocked Marshalmocked Sprintf")

	// inside mockSampleClz3.methodThatUsesMultileGlobalFunctions: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	sc3 := mockSampleClz3{}
	assert.True(sc3.methodThatUsesMultileGlobalFunctions("format", "value") == "\"format\"mocked Sprintf")
}

func TestClonedFuncs(t *testing.T) {
	assert := require.New(t)

	// setup function mocks
	jsonMock.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	fmtMock.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	assert.True(functionThatUsesGlobalFunction_clone("format", "value") == "mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	assert.True(functionThatUsesMultileGlobalFunctions_clone("format", "value") == "mocked Marshalmocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions2: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	assert.True(functionThatUsesMultileGlobalFunctions2_clone("format", "value") == "\"format\"mocked Sprintf")

}

package yaml

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSampleClz(t *testing.T) {
	assert := require.New(t)

	// inside mockSampleClz.methodThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	sc := mockSampleClz{}
	sc.mock_mockSampleClz_methodThatUsesGlobalFunction_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	assert.True(sc.methodThatUsesGlobalFunction("format", "value") == "mocked Sprintf")

	// inside mockSampleClz2.methodThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	sc2 := mockSampleClz2{}
	sc2.mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	sc2.mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_json.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	assert.True(sc2.methodThatUsesMultileGlobalFunctions("format", "value") == "mocked Marshalmocked Sprintf")

	// inside mockSampleClz3.methodThatUsesMultileGlobalFunctions: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	sc3 := mockSampleClz3{}
	sc3.mock_mockSampleClz3_methodThatUsesMultileGlobalFunctions_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	assert.True(sc3.methodThatUsesMultileGlobalFunctions("format", "value") == "\"format\"mocked Sprintf")
}

func TestClonedFuncs(t *testing.T) {
	assert := require.New(t)

	c := &mockFmtclonedFuncs{}
	c.mock_mockFmtclonedFuncs_functionThatUsesGlobalFunction_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	c.mock_mockFmtclonedFuncs_functionThatUsesMultileGlobalFunctions_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	c.mock_mockFmtclonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	c.mock_mockFmtclonedFuncs_functionThatUsesMultileGlobalFunctions_json.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)

	// inside functionThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	assert.True(c.functionThatUsesGlobalFunction("format", "value") == "mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	assert.True(c.functionThatUsesMultileGlobalFunctions("format", "value") == "mocked Marshalmocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions2: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	assert.True(c.functionThatUsesMultileGlobalFunctions2("format", "value") == "\"format\"mocked Sprintf")
}

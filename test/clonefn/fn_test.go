package clonefn

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClonedFuncs(t *testing.T) {
	assert := require.New(t)
	c := &clonedFuncs{}

	// setup function mocks
	c.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_json.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	c.mock_clonedFuncs_functionThatUsesGlobalFunction_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	c.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")
	c.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	assert.True(c.functionThatUsesGlobalFunction("format", "value") == "mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	assert.True(c.functionThatUsesMultileGlobalFunctions("format", "value") == "mocked Marshalmocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions2: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	assert.True(c.functionThatUsesMultileGlobalFunctions2("format", "value") == "\"format\"mocked Sprintf")

}

func Test_functionThatUsesFunctionFromSameRoot(t *testing.T) {
	assert := require.New(t)

	c := &mockCallee{}
	c.mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo.On("Dummy").Return("Mocked Dummy")
	s := c.functionThatUsesFunctionFromSameRoot()
	assert.Equal("Mocked Dummy", s)
}

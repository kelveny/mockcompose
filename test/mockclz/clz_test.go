package mockclz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test if we can clone from source class correctly
func TestClonedClz(t *testing.T) {
	assert := require.New(t)

	o := cloneSourceClz{}
	s, err := o.Unnamed("string", 1, nil)
	assert.True(s == "")
	assert.True(err.Error() == "Invalid arguments")

	s, err = o.Unnamed2("string", 1, nil)
	assert.True(s == "")
	assert.True(err.Error() == "Invalid arguments")

	s, err = o.Unnamed("string", 1, make(chan<- string))
	assert.True(s == "s: string, n: 1")
	assert.True(err == nil)

	s = o.Variadic("%s=%s", "key", "value")
	assert.True(s == "key=value")

	s = o.Variadic2("%s=%d", "value", 1)
	assert.True(s == "value=1")

	s = o.Variadic3("value1", "value2")
	assert.True(s == "value1, value2")

	s = o.Variadic4("value1", 1)
	assert.True(s == "\"value1\", 1")

	s = o.CollapsedParams([]byte("value1"), []byte("value2"))
	assert.True(s == "value1value2")

	x, y, s := o.CollapsedReturns()
	assert.True(x == 1)
	assert.True(y == 2)
	assert.True(s == "")
}

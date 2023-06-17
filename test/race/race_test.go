package race

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

// go test --race will fail the test
// testify is primarily used for unit testing, it
// does not specifically focus on data race testing
func TestRace(t *testing.T) {
	r := &raceMock{}

	r.On("WorkRun", mock.Anything).Return().Run(func(args mock.Arguments) {
		wg := args.Get(0).(*sync.WaitGroup)
		wg.Done()
	})

	r.RaceRun(10)

	r.AssertNumberOfCalls(t, "WorkRun", 10)
}

package race2

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

// go test --race will succeed in this test
//
// Go context is thread-safe immutable object, it is safe to mock
// functions that have immutable thread-safe parameters
func TestRace(t *testing.T) {
	r := &raceMock{}

	r.On("WorkRun", mock.Anything).Return().Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)

		wg := ctx.Value("wg").(*sync.WaitGroup)
		wg.Done()
	})

	r.RaceRun(10)

	r.AssertNumberOfCalls(t, "WorkRun", 10)
}

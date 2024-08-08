package race

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

// go test --race will fail this test
//
// It is not a good practice to mock functions that contain parameters having
// protected fields for multi-threaded access. The subject parameter is usually
// a struct object that has fields be protected by a mutex and may be accessed
// from other threads
//
// Inside testify Mock.Called() implementation, it performs read on the passed in
// parameeter wihout any synchronization on the parameter (it has no idea about the
// dynamic detail inside the object), this can cause "go test -race" to report
// data race errors
//
// A possible solution is to use context.Context, see race2 example

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

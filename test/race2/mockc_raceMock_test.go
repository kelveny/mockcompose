// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package race2

import (
	"context"
	"fmt"
	"sync"

	"github.com/stretchr/testify/mock"
)

type raceMock struct {
	race
	mock.Mock
}

func (m *raceMock) WorkRun(ctx context.Context) {

	m.Called(ctx)

}

func (r *raceMock) RaceRun(runners int) {
	wg := &sync.WaitGroup{}
	wg.Add(runners)
	ctx := context.WithValue(context.Background(), "wg", wg)
	for i := 0; i < runners; i++ {
		fmt.Printf("raceRun start runner %d\n", i)
		go r.WorkRun(ctx)
	}
	wg.Wait()
	fmt.Printf("raceRun done\n")
}

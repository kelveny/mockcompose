// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package race

import (
	"fmt"
	"sync"

	"github.com/stretchr/testify/mock"
)

type raceMock struct {
	race
	mock.Mock
}

func (m *raceMock) WorkRun(wg *sync.WaitGroup) {

	m.Called(wg)

}

func (r *raceMock) RaceRun(runners int) {
	wg := sync.WaitGroup{}
	wg.Add(runners)
	for i := 0; i < runners; i++ {
		fmt.Printf("raceRun start runner %d\n", i)
		go r.WorkRun(&wg)
	}
	wg.Wait()
	fmt.Printf("raceRun done\n")
}

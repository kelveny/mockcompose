package race

import (
	"fmt"
	"sync"
)

type Race interface {
	WorkRun(wg *sync.WaitGroup)

	RaceRun(runners int)
}

type race struct {
	lock  sync.Mutex
	count int
}

func (r *race) WorkRun(wg *sync.WaitGroup) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.count += 1
	wg.Done()
}

func (r *race) RaceRun(runners int) {
	wg := sync.WaitGroup{}
	wg.Add(runners)
	for i := 0; i < runners; i++ {
		fmt.Printf("raceRun start runner %d\n", i)
		go r.WorkRun(&wg)
	}
	wg.Wait()

	fmt.Printf("raceRun done\n")
}

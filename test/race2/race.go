package race2

import (
	"context"
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

func (r *race) WorkRun(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.count += 1

	wg := ctx.Value("wg").(*sync.WaitGroup)
	wg.Done()
}

func (r *race) RaceRun(runners int) {
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

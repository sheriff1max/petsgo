package main

import (
	"fmt"
	"sync"
)


var wg_race sync.WaitGroup

type Counter struct {
	counter int
	sync.Mutex
}

func (obj *Counter) Inc() {
	obj.Lock()
	obj.counter++
	obj.Unlock()
}

func (obj *Counter) Value() int {
	return obj.counter
}

func incRacy(counter *int, n int) {
	wg_race.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			(*counter)++
			wg_race.Done()
		}()
	}
}

func main() {
	n := 1000

	// 1. Data Race
	counter_race := 0
	incRacy(&counter_race, n)
	fmt.Printf("counter_race = %d\n", counter_race)

	// 2. With Mutex
	c := &Counter{}
	var wg sync.WaitGroup
	for range n {
		wg.Go(
			func() {
				c.Inc()
			},
		)
	}
	wg.Wait()
	fmt.Println(c.Value() == n)
}
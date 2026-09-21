package main

import (
	"fmt"
	"time"
	"sync"
)


var wg sync.WaitGroup


func process(n, k int) []int {
	ch := make(chan struct{}, k)
	res := make([]int, n)

	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			ch <- struct{}{}
			defer func(){ <- ch }()

			time.Sleep(5 * time.Millisecond)
			res[idx] = idx
		}(i)
	}
	wg.Wait()
	return res
}

func main() {
	t := time.Now()
	res := process(10, 3)
	fmt.Println("Прошло время = ", time.Since(t))

	s := 0
	for _, v := range res {
		s += v
	}
	fmt.Println(len(res) == 10 && s == 45 && res[9] == 9)
}
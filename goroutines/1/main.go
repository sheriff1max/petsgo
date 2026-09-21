package main

import (
	"fmt"
	"sync"
)


var wg sync.WaitGroup


func parallelSquares(n int) []int {
	arr := make([]int, n)

	squar := func(n int) {
		arr[n] = n * n
		wg.Done()
	}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go squar(i)
	}
	wg.Wait()
	return arr
}

func main() {
	res := parallelSquares(10)
	s := 0
	for _, v := range res {
		s += v
	}
	fmt.Println(len(res) == 10 && res[0] == 0 && res[9] == 81 && s == 285)
}
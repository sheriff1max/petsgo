package main

import (
	"fmt"
)


func generate(n int) <- chan int {
	ch := make(chan int)

	go func() {
		for i := 1; i <= n; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func main() {
	s, n := 0, 0
	for v := range generate(10) {
		s += v
	}
	for range generate(5) {
		n++
	}
	fmt.Println(s == 55 && n == 5)
}
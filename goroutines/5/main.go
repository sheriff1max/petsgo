package main

import (
	"fmt"
	"time"
)


func delayed(v int, d time.Duration) <- chan int {
	ch := make(chan int)

	go func() {
		time.Sleep(d)
		ch <- v
		close(ch)
	}()
	return ch
}

func fetchWithTimeout(ch <- chan int, timeout time.Duration) (int, bool) {


	select {
	case v, ok := <- ch:
		if !ok {
			return 0, false
		}
		return v, true
	case <- time.After(timeout):
		return 0, false
	}
}

func main() {
	v, ok1 := fetchWithTimeout(delayed(42, 10*time.Millisecond), 100*time.Millisecond)
	_, ok2 := fetchWithTimeout(make(chan int), 50*time.Millisecond)
	fmt.Println(ok1 && v == 42 && !ok2)
}

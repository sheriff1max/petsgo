package main

import (
	"fmt"
)


func gen(n int, out chan <- int) {
	for i := 1; i <= n; i++ {
		out <- i
	}
	close(out)
}

func square(in <- chan int, out chan <- int) {
	for val := range in {
		out <- val * val
	}
	close(out)
}

func filterEven(in <- chan int, out chan <- int) {
	defer close(out)
	for val := range in {
		if val % 2 == 0 {
			out <- val
		}
	}
}

func collect(in <- chan int) []int {
	var res []int
	for val := range in {
		res = append(res, val)
	}
	return res
}

func collectPipeline(n int) []int {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go gen(n, ch1)
	go square(ch1, ch2)
	go filterEven(ch2, ch3)
	return collect(ch3)
}

func main() {
	res := collectPipeline(10)
	s := 0
	for _, v := range res {
		s += v
	}
	fmt.Println(len(res) == 5 && s == 220 && res[0] == 4 && res[4] == 100)
}
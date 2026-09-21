package main

import (
	"fmt"
)


func producer(out chan <- int, n int) {
	defer close(out)
	for i := 1; i <= n; i++ {
		out <- i
	}
}

func doubler(in <- chan int, out chan <- int) {
	defer close(out)
	for num := range in {
		out <- num * 2
	}
}

func collector(in <- chan int) []int {
	var res []int
	for num := range in {
		res = append(res, num)
	}
	return res
}

func runPipeline(n int) []int {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go producer(ch1, n)
	go doubler(ch1, ch2)
	return collector(ch2)
}

func main() {
	out := runPipeline(5) // collector после doubler после producer
	fmt.Println(len(out) == 5 && out[0] == 2 && out[4] == 10)
}

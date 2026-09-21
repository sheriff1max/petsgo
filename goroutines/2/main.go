package main

import (
	"fmt"
)


func sumPart(left, right int, ch chan int) {
	sum := 0
	for i := left; i <= right; i++ {
		sum += i
	}
	ch <- sum
}

func sumAsync(n, workers int) int {
	ch := make(chan int, workers)

	one_part := n / workers

	for i := 0; i < workers - 1; i++ {
		sumPart(i * one_part + 1, (i + 1) * one_part, ch)
	}
	sumPart((workers - 1) * one_part + 1, n, ch)

	res := 0
	for i := 0; i < workers; i++ {
		res += <- ch
	}

	return res
}

func main() {
	fmt.Println(sumAsync(100, 4) == 5050 && sumAsync(10, 2) == 55 && sumAsync(7, 7) == 28)
}

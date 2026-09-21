package main

import (
	"fmt"
	"time"
	"math"
)


type offer struct {
	price, idx int
}

func cheapest(prices []int, delays []time.Duration) (price int, idx int) {
	ch := make(chan offer)

	for i := range len(prices) {
		go func(idx int) {
			time.Sleep(delays[idx])
			ch <- offer{price: prices[idx], idx: idx}
		}(i)
	}

	min_price, min_idx := math.MaxInt, 0
	for range len(prices) {
		val := <- ch
		if val.price < min_price || (val.price == min_price && val.idx < min_idx) {
			min_price = val.price
			min_idx = val.idx
		}
	}
	return min_price, min_idx
}

func main() {
	p1, i1 := cheapest([]int{100, 80, 90}, []time.Duration{30 * time.Millisecond, 10 * time.Millisecond, 20 * time.Millisecond})
	p2, i2 := cheapest([]int{50, 50, 60}, []time.Duration{10 * time.Millisecond, 10 * time.Millisecond, 10 * time.Millisecond})
	fmt.Println(p1 == 80 && i1 == 1 && p2 == 50 && i2 == 0)
}

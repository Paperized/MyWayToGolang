package main

import (
	"fmt"
	"signals/signals"
	"time"
)

func main() {
	curr := time.Now()
	totalSleptSecondsSg := signals.MakeSignal(0)
	totalSleptSecondsSg.Listen(func(i int, bs *signals.BaseSignal[int]) {
		fmt.Printf("Total slept seconds: %d\n", i)
	})

	sleep := signals.MakeSignal(0)
	sleep.ListenAsync(func(i int, bs *signals.BaseSignal[int]) {
		time.Sleep(time.Second * time.Duration(i))
		totalSleptSecondsSg.SetFromValue(func(x int) int { return x + i })

		// try this, you will see an unexpected error due to the async operations
		// totalSleptSecondsSg.Set(totalSleptSecondsSg.Get() + i)
	})

	sleep.ListenAsync(func(i int, bs *signals.BaseSignal[int]) {
		time.Sleep(time.Second * time.Duration(i))
		totalSleptSecondsSg.SetFromValue(func(x int) int { return x + i })

		// try this, you will see an unexpected error due to the async operations
		// totalSleptSecondsSg.Set(totalSleptSecondsSg.Get() + i)
	})

	sleep.ListenAsync(func(i int, bs *signals.BaseSignal[int]) {
		time.Sleep(time.Second * time.Duration(i))
		totalSleptSecondsSg.SetFromValue(func(x int) int { return x + i })

		// try this, you will see an unexpected error due to the async operations
		// totalSleptSecondsSg.Set(totalSleptSecondsSg.Get() + i)
	})

	sleep.Set(3)
	fmt.Printf("Time elapsed: %v", time.Since(curr))
}

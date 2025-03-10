package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ch1, ch2 := make(chan int), make(chan int)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(2 * time.Millisecond)
		ch2 <- 2
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case val := <-ch1:
		fmt.Println("ch1", val)
	case val := <-ch2:
		fmt.Println("ch2", val)
	case <-time.After(5 * time.Millisecond):
		fmt.Println("timeout ")
	}
}

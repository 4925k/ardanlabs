package main

import (
	"fmt"
	"time"
)

func main() {
	// GO ROUTINES
	go fmt.Println("goroutine")
	fmt.Println("main")

	for i := range 3 {
		go func() {
			fmt.Println("normal", i)
		}()

		// passing parameters
		go func(int) {
			fmt.Println("parameter", i)
		}(i)

		i := i // shadowing
		go func() {
			fmt.Println("shadowing", i)
		}()
	}

	time.Sleep(1 * time.Second)

	// -----------------------------------------------
	// CHANNELS

	ch := make(chan string)
	// ch <- "hi" // this will deadlock
	go func() {
		ch <- "hi" // sending to channel
	}()
	msg := <-ch // receiving from channel
	fmt.Println(msg)

	go func() {
		for i := range 3 {
			ch <- fmt.Sprintf("hi %d", i)
		}
		close(ch) // notify that channel is closed
	}()

	for msg := range ch {
		fmt.Println(msg)
	}

	// // for range on channels do this
	// for {
	// 	msg, ok := <-ch
	// 	if !ok {
	// 		break
	// 	}

	// 	fmt.Println("got: ", msg)
	// }

	msg = <-ch // ch is closed. you will get nothing
	fmt.Println("message from closed channel:", msg)

	// how to know if the channel is closed?
	msg, ok := <-ch
	fmt.Printf("msg: %s\tok: %t\n", msg, ok)

	// sending on closed channels
	// ch <- "hi" // PANIC: sending on closed channel
}

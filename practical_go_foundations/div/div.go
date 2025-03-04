package main

import (
	"fmt"
	"log"
)

func main() {
	val, err := safeDiv(10, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("val: %d\n", val)
}

func safeDiv(a, b int) (ret int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("error: %v", e)
		}
	}()

	return a / b, nil
}

func div(a, b int) int {
	return a / b
}

package main

import (
	"fmt"
	"sync"
	"time"
)

type Payment struct {
	From   string
	To     string
	Amount float64

	once sync.Once
}

func main() {
	p := Payment{
		From:   "Alice",
		To:     "Bob",
		Amount: 100,
	}

	p.Process()
	p.Process()
}

func (p *Payment) Process() {
	t := time.Now()
	p.once.Do(func() { p.process(t) })

	// p.once.Do(p.process(t)) // how to pass arguments

}

func (p *Payment) process(t time.Time) {
	ts := t.Format(time.RFC3339)
	fmt.Printf("[%s] from: %s to: %s amount: %.2f\n", ts, p.From, p.To, p.Amount)

}

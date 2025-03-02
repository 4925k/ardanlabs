package main

import (
	"fmt"
	"log"
)

const (
	maxX = 1000
	maxY = 1000
)

// Item is an item in the game
type Item struct {
	X int
	Y int
}

type Player struct {
	Name string
	Item // Embed
}

func main() {
	var i1 Item
	fmt.Println(i1)
	fmt.Printf("%#v\n", i1)

	i2 := Item{1, 2}
	fmt.Printf("i2: %#v\n", i2)

	i3 := Item{X: 1, Y: 2}
	fmt.Printf("i3: %#v\n", i3)

	i4, err := NewItem(10, 10)
	if err != nil {
		log.Fatalf("creating item: %v", err)
	}
	fmt.Printf("i4: %#v\n", i4)

	i4.Move(20, 20)
	fmt.Printf("i4: %#v\n", i4)

	p1 := Player{
		Name: "p1",
		Item: Item{500, 500},
	}
	fmt.Printf("p1: %#v\n", p1)
	fmt.Printf("p1.Item.X %#v\n", p1.Item.X)

	ms := []mover{&i2, &p1, i4} // cant use i1 here because it is a value.
	moveAll(ms, 0, 0)

}

func moveAll(ms []mover, x, y int) {
	for _, m := range ms {
		m.Move(x, y)
	}
}

type mover interface {
	Move(x, y int)
	// Move(int, int)
}

// NewItem returns a pointer to an Item
// func NewItem(x, y int) Item
// func NewItem(x, y int) *Item
// func NewItem(x, y int) (Item, error)
func NewItem(x, y int) (*Item, error) { // best practices
	if x < 0 || y < 0 || x > maxX || y > maxY {
		return nil, fmt.Errorf("invalid position: %d, %d", x, y)
	}

	// i is allocated on the heap. scope within the function
	i := Item{
		X: x,
		Y: y,
	}

	// the go compile does escape analysis and will allocate i on the heap
	// because i is going to outlive the function.
	return &i, nil
}

// i is a value receiver
func (i *Item) Move(x, y int) {
	i.X = x
	i.Y = y
}

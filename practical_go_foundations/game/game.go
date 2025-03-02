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
	Keys []Key
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

	k := Jade
	fmt.Printf("k: %s\n", k)

	err = p1.FoundKey(k)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("p1.Key: %#v\n", p1.Keys)

	err = p1.FoundKey(invalidKey)
	if err != nil {
		log.Fatal(err)
	}

	// time.Time imports json.Marshaler interface
	// json.NewEncoder(os.Stdout).Encode(time.Now())
}

const (
	Jade Key = iota + 1
	Copper
	Crystal
	invalidKey // internal
)

// Implement fmt.Stringer interface
func (k Key) String() string {
	switch k {
	case Jade:
		return "Jade"
	case Copper:
		return "Copper"
	case Crystal:
		return "Crystal"
	}

	return fmt.Sprintf("unknown key: %d", k)
}

// FoundKey will add k to Key if its not there
// error if k is not one of the known keys
func (p *Player) FoundKey(k Key) error {
	if k < Jade || k >= invalidKey {
		return fmt.Errorf("unknown key: %d", k)
	}

	if !p.containsKey(k) {
		p.Keys = append(p.Keys, k)
	}

	return nil
}

func (p *Player) containsKey(k Key) bool {
	for i := range p.Keys {
		if p.Keys[i] == k {
			return true
		}
	}

	return false
}

type Key byte

func moveAll(ms []mover, x, y int) {
	for _, m := range ms {
		m.Move(x, y)
	}
}

type mover interface {
	Move(x, y int)
	// Move(int, int)
}

/*
NewItem returns a pointer to an Item
func NewItem(x, y int) Item
func NewItem(x, y int) *Item
func NewItem(x, y int) (Item, error)
The following signature is the best practice
*/
func NewItem(x, y int) (*Item, error) {
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

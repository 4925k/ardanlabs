package main

import "fmt"

func main() {
	banner("Go", 6)
	banner("G©", 6)

	s := "G©"
	fmt.Println("length of G©: ", len("G©"))

	// code point = rune ~= unicode character
	for i, r := range "G©" {
		fmt.Println(i, r)
		fmt.Printf("%c of type %T\n", r, r)
	}

	// byte (uint8)
	// rune (int32)

	b := s[1]
	fmt.Printf("%c of type %T\n", b, b)

	x, y := 1, "1"
	fmt.Printf("x=%#v y=%#v\n", x, y) // use %#v for debugging

	fmt.Println(isPalindrome("Go"))
	fmt.Println(isPalindrome("racecar"))
	fmt.Println(isPalindrome("G©G"))
}

func isPalindrome(s string) bool {
	rs := []rune(s) // get slice of runes out of s -> making it unicode friendly
	for i := 0; i < len(rs)/2; i++ {
		if rs[i] != rs[len(rs)-1-i] {
			return false
		}
	}

	return true
}

func banner(text string, width int) {
	padding := (width - len(text)) / 2

	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}

	fmt.Println(text)

	for i := 0; i < width; i++ {
		fmt.Print("-")
	}

	fmt.Println()

}

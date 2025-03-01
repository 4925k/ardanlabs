package main

import (
	"fmt"
	"sort"
)

func main() {
	var s []int
	fmt.Printf("len: %d\n", len(s))

	if s == nil {
		fmt.Println("nil slice")
	}

	s2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("s2 = %#v\n", s2)

	// SLICING OPERATION
	s3 := s2[2:6]
	fmt.Printf("s3 = %#v\n", s3)

	// fmt.Println(s2[:100]) // PANIC

	s3 = append(s3, 100)
	fmt.Printf("s3 = %#v\n", s3)
	fmt.Printf("s2 = %#v\n", s2) // s2 is changed as well
	fmt.Printf("s2: len: %d cap: %d\n", len(s2), cap(s2))
	fmt.Printf("s3: len: %d cap: %d\n", len(s3), cap(s3))

	var s4 []int
	for i := 0; i < 1_000; i++ {
		s4 = appendInt(s4, i)
	}
	// fmt.Println("s4:", s4)
	// fmt.Println("s4:", s4[:2000]) // PANIC

	fmt.Println(concat([]string{"A", "B"}, []string{"C", "D"}))

	vs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(median(vs))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// copy in order to avoid changing the original array
	nums := make([]float64, len(values))
	copy(nums, values)

	sort.Float64s(nums) // NOTE this will sort the underlying array
	i := len(nums) / 2

	// if len(values)&1 == 1 { // binary calculation
	if len(nums)%2 == 1 {
		return nums[i]
	}

	return (nums[i-1] + nums[i]) / 2
}

func concat(s1, s2 []string) []string {
	// return append(s1, s2...) // my solution

	s := make([]string, len(s1)+len(s2))
	copy(s, s1)
	copy(s[len(s1):], s2)
	return s
}

// appendInt appends an int to a slice
// it shows the underlying mechanism of append
func appendInt(s []int, v int) []int {
	i := len(s)

	if len(s) < cap(s) { //enough space in the underlying array
		s = s[:len(s)+1]
	} else { // not enough capacity, need to grow
		fmt.Printf("grow: len: %d -> %d\n", len(s), 2*len(s)+1)
		s2 := make([]int, 2*len(s)+1)
		copy(s2, s)
		s = s2[:len(s)+1]
	}

	s[i] = v
	return s
}

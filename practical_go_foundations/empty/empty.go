package main

import "fmt"

func main() {
	var i any

	i = 8
	fmt.Printf("type: %T\tvalue:%v\n", i, i)

	i = "hello"
	fmt.Printf("type: %T\tvalue:%v\n", i, i)

	s := i.(string) // assertion
	fmt.Printf("type: %T\tvalue:%v\n", s, s)

	// t := i.(int) // will panic
	t, ok := i.(int)
	if !ok {
		fmt.Println("i is not an int")
	} else {
		fmt.Printf("type: %T\tvalue:%v\n", t, t)
	}

	switch i.(type) {
	case int:
		fmt.Println("i is an int")
	case string:
		fmt.Println("i is an string")
	default:
		fmt.Println("i is something else")
	}

	// Rule of thumb: Dont use any

	// fmt.Println("max int:", maxInts([]int{1, 2, 3, 4, 5}))
	// fmt.Println("max float64:", maxFloat64s([]float64{1.1, 2.2, 3.3, 4.4, 5.5}))

	// using generics
	fmt.Println("max int:", max([]int{1, 2, 3, 4, 5}))
	fmt.Println("max float64:", max([]float64{1.1, 2.2, 3.3, 4.4, 5.5}))

}

type Numbers interface {
	int | float64
}

// func max[T int | float64](nums []T) T {
func max[T Numbers](nums []T) T {

	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	return max

}

// maxInts returns the maximum value in the slice
func maxInts(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	return max
}

// maxFloat64s returns the maximum value in the slice
func maxFloat64s(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	return max
}

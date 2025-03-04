package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// WHAT IS THE MOST COMMON WORD IN THE FILE
	f, err := os.Open("sherlock.txt")
	if err != nil {
		log.Fatalf("reading file: %v", err)
	}
	defer f.Close()

	word, err := mostCommon(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("most common word: %s\n", word)
}

var wordRe = regexp.MustCompile(`\w+`)

func mostCommon(r io.Reader) (string, error) {
	freq, err := wordFrequency(r)
	if err != nil {
		return "", err
	}

	return maxWord(freq)
}

func wordFrequency(r io.Reader) (map[string]int, error) {
	s := bufio.NewScanner(r)

	count := make(map[string]int)
	// reading file line by line
	for s.Scan() {
		word := wordRe.FindAllString(s.Text(), -1)
		for _, w := range word {
			count[strings.ToLower(w)]++
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	return count, nil
}

func maxWord(freq map[string]int) (string, error) {
	if len(freq) == 0 {
		return "", fmt.Errorf("empty map")
	}

	maxCount, maxWord := 0, ""
	for w, c := range freq {
		if c > maxCount {
			maxCount = c
			maxWord = w
		}
	}

	return maxWord, nil
}

func mapDemo() {
	stocks := map[string]float64{
		"GOOG": 1000,
		"MSFT": 2000,
		"FB":   3000,
	}

	symbol := "MSFT"

	// zero vs nil values in maps
	if price, ok := stocks[symbol]; ok {
		fmt.Printf("price of %s is %.2f\n", symbol, price)
	} else {
		fmt.Printf("no price for %s\n", symbol)
	}

	// looping
	for k := range stocks {
		fmt.Printf("stock: %s\t price: %.2f\n", k, stocks[k])
	}

	// deleting
	delete(stocks, "MSFT")
	delete(stocks, "MSFT") // does nothing
}

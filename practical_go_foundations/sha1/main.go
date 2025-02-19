package main

import (
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	sha, err := sha1Sum("testfile.gz")
	if err != nil {
		log.Fatalf("generating sha1: %v", err)
	}

	fmt.Printf("sha1sum: %s\n", sha)

	sha, err = sha1Sum("testfile")
	if err != nil {
		log.Fatalf("generating sha1: %v", err)
	}

	fmt.Printf("sha1sum: %s\n", sha)
}

func sha1Sum(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var r io.Reader = file

	if strings.HasSuffix(fileName, ".gz") {
		gz, err := gzip.NewReader(file)
		if err != nil {
			return "", err
		}
		defer gz.Close()

		r = gz
	}

	w := sha1.New()

	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", w.Sum(nil)), nil
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Defining a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", true, "Count bytes")
	// Parsing the flags provided by the user
	flag.Parse()

	// Calling the count function to counter the number of words (or lines)
	// received from the standard input and printing it out
	fmt.Println(count(os.Stdin, *lines, *bytes))
}

func count(r io.Reader, clines bool, cbytes bool) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// Defining a counter for word or byte count
	c := 0

	switch {
	case clines:
		scanner.Split(bufio.ScanLines)
	case cbytes:
		scanner.Split(bufio.ScanBytes)
	default:
		scanner.Split(bufio.ScanWords)
	}

	// For every bytes scanned, increment the counter
	for scanner.Scan() {
		c++
	}

	return c
}

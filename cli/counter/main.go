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

func count(r io.Reader, clines bool, cbytes bool) (int, int) {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// Defining a counter for word count and byte count
	wc := 0
	bc := 0

	// If the clines flag is not set, we want to count words so we set
	// the scanner split types to words (default is split by lines)
	if !clines {
		scanner.Split(bufio.ScanWords)

		// For every word scanned, increment the counter
		for scanner.Scan() {
			wc++
		}
	}

	if cbytes {
		scanner.Split(bufio.ScanBytes)

		// For every bytes scanned, increment the counter
		for scanner.Scan() {
			bc++
		}
	}

	return wc, bc
}

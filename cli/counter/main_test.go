package main

import (
	"bytes"
	"testing"
)

// TestCountWordsAndBytes tests the count function set to count words and no byte count
func TestCountWordsWithBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\n")

	expWC := 3
	expBC := 0
	expBC2 := len(b.String())

	resWC, resBC := count(b, false, false)
	_, resBC2 := count(b, false, true)

	if resWC != expWC {
		t.Errorf("Expected %d, got %d instead.\n", expWC, resWC)
	}

	if resBC != expBC {
		t.Errorf("Expected %d, got %d instead.\n", expBC, resBC)
	}

	if resBC2 != expBC2 {
		t.Errorf("Expected %d, got %d instead \n", expBC2, resBC2)
	}
}

// TestCountLines tests the count function set to count lines
// func TestCountLines(t *testing.T) {
// 	b := bytes.NewBufferString("line1\nline2")

// 	exp := 2

// 	res, _ := count(b, true, true)

// 	if res != exp {
// 		t.Errorf("Expected %d, got %d instead.\n", exp, res)
// 	}
// }

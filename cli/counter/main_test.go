package main

import (
	"bytes"
	"testing"
)

// TestCountWordsAndBytes tests the count function set to count words and no byte count
func TestCountWordsWithBytes(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\n")

	expected := 3

	result := count(b, false, false)

	if expected != result {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}

// TestCountLines tests the count function set to count lines
func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("line1\nline2")

	expected := 2

	result := count(b, true, false)

	if result != expected {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}

// TestCountBytes tests the count function set to count bytes
func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("line1\nline2")

	expected := b.Len()

	result := count(b, false, true)

	if result != expected {
		t.Errorf("Expected %d, got %d instead.\n", expected, result)
	}
}

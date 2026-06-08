package greetings

import (
	"regexp"
	"testing"
)


func TestHelloName(t *testing.T) {
	name := "josh"
	want := regexp.MustCompile(`\b`+name+`\b`)

	msg, err := Hello(name)

	if !want.MatchString(msg) || err != nil {
		t.Errorf(`Hello("josh") = %q, %v, want match for %#q`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")

	if err == nil {
		t.Errorf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
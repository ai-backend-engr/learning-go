package greetings

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

// Hello returns a greeting for the named user
func Hello(name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}
	
	message := fmt.Sprintf(randomGreeting(), name)
	return message, nil
}

func Hellos(names []string) (map[string]string, error) {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	var greetings = make(map[string]string)

	for _, name := range names {
		message, err := Hello(name)
		if err != nil {
			log.Fatal(err)
		}

		greetings[name] = message
	}

	return greetings, nil
}

// generates random greetings
func randomGreeting() string {
	formats := []string{
		"Hi, %v. Welcome!",
        "Great to see you, %v !",
        "Hail, %v ! Well met!",
	}

	return formats[rand.Intn(len(formats))]
}
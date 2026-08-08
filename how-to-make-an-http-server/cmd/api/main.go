package main

import (
	"log"

	"example.com/htmahs/internal/app"
)

func main() {
	app, err := app.New()

	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}

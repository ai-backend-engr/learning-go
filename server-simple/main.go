package main

import (
	"fmt"
	"log"
	"net/http"
)

// Handler function to handle incoming request
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Congratulations you have connected to the server home.\n")
}

func main() {
	// 1. Register a route and its handler func
	http.HandleFunc("/", homeHandler)

	// 2. Define the port number
	port := ":8080"
	fmt.Printf("Starting http server on port %s...\n", port)

	// 3. Start the listener and server loop
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start : %v", err)
	}
}

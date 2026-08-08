package main

import (
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s request", r.URL)
	io.WriteString(w, "Hello from root")
}

func main() {
	http.HandleFunc("/", getRoot)

	http.ListenAndServe(":3000", nil)
}

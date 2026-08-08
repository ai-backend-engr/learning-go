// STEP 4
package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			count, err := strconv.Atoi(r.FormValue("counter"))
			if err != nil {
				fmt.Fprintf(w, "<p>Input a valid number</p>")
			}

			count++

			fmt.Fprintf(w, `
				<body>
					<form action="/" method="POST">
						<label>Counter</label>
						<input name="counter" value="%d" readonly>
						<button type="submit">Add</button>
					</form>
					<a href="/">Reset</a>
				</body>
			`, count)

			return
		}

		fmt.Fprintf(w, `
            <body>
                <form action="/" method="POST">
                    <label>Counter</label>
                    <input name="counter" value="1" readonly>
                    <button type="submit">Add</button>
                </form>
            </body>
            <a href="/">Reset</a>
        `)
	})

	http.ListenAndServe(":8080", nil)
}

// STEP 3
/*
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("q") != "" {
			fmt.Fprintf(w, "<h1>Hello %s</h1>", r.FormValue("q"))
		}

		fmt.Fprintf(w, `
			<body>
                <form action="/" method="GET">
                    <label>Enter your name</label>
                    <input name="q">
                    <button type="submit">Submit</button>
                </form>
            </body>
		`)
	})

	http.ListenAndServe(":8080", nil)
}
*/

// STEP 2
/*
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Hello, World!</h1>")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

// STEP 1
/*
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
*/

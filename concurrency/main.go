package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	numChan := make(chan int)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go printNumbers(i, &wg, numChan)
	}

	generateNumbers(5, numChan)

	close(numChan)

	fmt.Println("Waiting for goroutines to complete")
	wg.Wait()
	fmt.Println("Done!")
}

func printNumbers(idx int, wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()

	for num := range ch {
		fmt.Printf("%d: read %d from channel\n", idx, num)
	}
}

func generateNumbers(total int, ch chan<- int) {
	for idx := 1; idx <= total; idx++ {
		fmt.Printf("Generating number %d\n", idx)
		ch <- idx
	}
}

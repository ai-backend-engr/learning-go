package main

import (
	"context"
	"fmt"
	"time"
)

func producer(ctx context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	printCh := make(chan int)

	go worker(ctx, printCh)

	for i := 1; i <= 3; i++ {
		printCh <- i
	}

	cancelCtx()

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("producer finished!")
}

func worker(ctx context.Context, printCh chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("worker err: %s\n", err)
			}
			fmt.Println("worker finished!")
			return
		case num := <-printCh:
			fmt.Printf("doAnother: %d\n", num)
		}

	}
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "myKey", "myValue")
	producer(ctx)
}

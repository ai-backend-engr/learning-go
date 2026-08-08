# Concurrency in Go
GO's core feature includes running programs concurrently. As it adapts to the ability of modern computers to run multiple streams of code simultaneously. For programmers to write concurrent application in Go we use these 2 features, goroutine and channels. Goroutine helps us spin of concurrent code while channels helps us safely communicate between the concurrent code.

## Goroutines
In modern computers the CPU or processor has one or more cores each capable of running a stream of code simultaneously. For our program to take advantage of this speed we need to split our code into multiple stream which Go was designed to make easy.

Goroutine is what makes this possible. Goroutine is a special kind of function that allow another code to run while the still running the current function. Usually a function would run in the foreground and block operations until it is done, goroutines runs in the background allowing other code, functions to run simultenously with it.

One thing that even makes Goroutine a sweet spot is the fact that it can run concurrently on a single core and run in parallel on multiple cores. But you understand how the CPU runs tasks with its FIFO etc. This makes a multi core run tasks both simultaneously and in parallel, which makes Go codes scalable.

When two different operations do not depend on each other and each take a significant amount of time to return, it calls for the need of a goroutine, which theoritically cuts the turnaround time in half.

**NOTE:** 
Just adding the go keyword before the function would return nothing. 

Hence you would need a way for the program to know that both goroutines have finished. If you do not wait they would never run or run incompletely. We can use a synchronization primitive from sync package - `WaitGroup` to do that. A synchronization primitive are used to synchornize the various parts of an application, in our case that would be to keep track of when each goroutine is done running.

Reason being that, we didnt wait for the goroutine to return before returning from the main func. Hence, in this case we need a synchronization primitive such as a `WaitGroup`. Synchronization primitives are designed to synchronize various parts of our application. In our case it tells us if we are done with the goroutine functions.

***NOTE: Goroutines do not have return keywords like standard function do*** 

However you can still use goroutines to call a function that has a return keyword. Still you wont be able to access those return data, they would be thrown out. There is a feature which we can use to access them, ***channels***

## Channels
If you're not careful you would run into problems that exists only in concurrent applications. Issues like race conditions. 

Go creators knew you would run into this situations so, they built a feature called ***channels*** for you, a sort of pipe connecting 2 goroutines. It still doesnt mean you cannot still encounter the race conditions if you are not careful.

```go
// creating a channel
channel := make(chan int)

// after creating a channel
// you can send or receive data using the arrow-looking operator (<-)

// To write to a channel
channel <- 10

// To read from a channel
intVar := <- channel

// Looping a channel with range
for num := range channel {
    if num < 1 {
        break;
    }
}

// what if you want a function to perform only read or write functions
func writeChannel (ch chan <- int) {
    // ch is write only
}

func readChannel (ch <- chan int) {
    // ch is read only
}

func readWriteChannel (ch chan int) {
    // ch is read and write
}
```
When a channel is no longer being used it is good practice to close it with the built in `close()` function to prevent ***memory leak***. ***A memory leak is when a program makes use of a memory and doesnt return it back to the computer resulting in the program using up more memory than expected.***

While using `chan int` is possible for function signatures, it is good practice to avoid it, to prevent something known as ***deadlock***. This is when a program is waiting for a process to finish meanwhile the process is waiting for that program to finish as well.

This deadlock is possible because of a very important principle regarding how channels work behind the scene it is important to know that channels are blocking. What it means is that when you write to a channel it block waiting for a process to read from it before continuing that is also the case for write. So using `chan int` might result to situation where the code is expected to read but you are performing a write situations like this causes you to enter a deadlock because that process would never be written to.

I think I am in love because as we understand how goroutine works, we can scale our application by running mulitple goroutines to act on the same data.

Before you decide to turn your entire app into a goroutine. Wait and listen. Doing that could result to resource starvation. A process where the computer doesn't have enough resources to do its work. What do you mean? Here is what I mean, you spin off a thousand goroutines, and the computer has less cores to actually do this job, then you, also understand that the ogoroutine take some time to swtich context. Combined you can see that the computer might spend its time switching context instead of actually doing its job. Know what system you are running your go program on and define the goroutines, which is safer as you know the capacity of the server you are running on.
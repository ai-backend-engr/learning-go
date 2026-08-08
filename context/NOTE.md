# Context in Go
Some times you need to obtain information about the environment a function is running. Okay take this instance, we have a multiplexer doing some multiplexing so you would expect just a URL right, but what about the client info... Is this client still connected? We also need that info so we can manage our resources. Things like this is what context helps us navigate correctly.

## Creating context
`context.TODO()`: ***this can be useful when creating a starting context. You should use this when you are not sure which context to use***

`context.Background()`: ***this is also useful when creating a starting context. You should however use this when you want to start a known context***

## Adding data to context
`context.WithValue()`: ***this is used to add value to a context***
```go
// it accepts three values the context, key, and value
ctx := context.Background()
ctx = context.WithValue(ctx, "name", "Joshua")

// NOTE: adding value to a context creates a new context and the values added to contexts are immutable 
// it wraps the parent inside the new one and outputs the new one.
```

## Accessing data in context
`context.Value()`: ***this is used to access the value in a context by providing a key***
```go
context.Value("name")
// output => "Joshua"
```

Context are the not just there for you to pass data around usually you would pass data around using variables and parameters, there are only a few cases where it makes sense to pass data around using context. Maybe if you need to log a users information inside multiple function but the functions themselves do not need the users data.

## Ending a context
Another important tooling that context provides is its ability to signal to any function running a context that the context is done. This is important so that the function can immediately stop running and free up resources for other part of the code to run.

There are three methods in which we can use to end a context;
### 1. Determining if a context is Done.
How to figure out if a context has ended is through the `Done()` method provided by context. It returns a channel that is closed when the context is `Done()`

We can combine this periodic check of the `Done()` method, the processing of a functionality and the select method to do even more useful work, like writting to serveral other channels almost simultenously while watching for when the context is canceled.
```go
ctx := context.Background()
resultChan := make(chan *WorkResult)

for {
    select {
        case <- ctx.Done():
        // the context is over, stop processing
        return
    result := <-chan resultChan:
        // perform some operation
    }
}
```
Each time a select keyword is found in a function, Go stops running the function and checks the cases of the function. It is also important to understand that the order of preference is not specified if more than one can run simultaneously.

### 2. Canceling a context
Cancelling a context is the most straightforward way to end a context. 

`context.WithValue()`: ***this is used to cancel a context, it returns a ctx and a cancel function to cancel the returned function***

```go
parentCtx := context.Background()
childCtx := context.WithValue(parentCtx, "key", "value")

// creating a cancel function
newCtx, cancelFunc := context.WithCancel(childCtx)
cancelFunc() // this function cancels the [newCtx]
```

`context.WithDeadline()`: ***this ***



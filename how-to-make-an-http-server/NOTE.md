# CREATING A PRODUCTION SERVER
We would learn how to create a production level server in Go...

## The Anatomy of a GO HTTP Server
Clients send requests to our server who is listening for a TCP connection using a mechanism known as TCP Listener. When a connection is established requests are passed through and parsed with a HTTP parser which creates a request object that is handled by our router which figures our where to send the request to be handled by our request handler. Our hanlder processes the request and writes a response through a ResponseWriter and the response is serialized and sent back to the client as response.

```go
// a handler func
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s request came in", r.URL)
}

func main() {
    // registering a handler
    http.HandleFunc("/", myHandler)

    // A simple Go server
    http.ListenAndServe(":3000", nil)
}
```
The method `HandleFunc()` handles a func as a req handler. A handler is simply something capable of handling an HTTP request nothing more, nothing less.

`http.HandleFunc()` registers `myHandler()` as a route handler

`ResponseWriter` is an interface. Because we might want to write responses back in different forms. But it must implement this interface.
```go
type ResponseWriter interface {
    WriteHeader
    Write 
}
```

`Request` is a struct. Because all request have the same shape. It is a pointer because we want to reference the exact request even deep in a middleware, and not waste space replicating the req object in memory
```go
type Request struct {

}
```

Everything starts from here

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```
Interface is just a behaviour. So when we say an object(struct) implements an interface, we are saying that the struct has that behaviour. Interface is one of flexibility over rigidity which supports structural typing.

Anybody who has the `Handler` behaviour can be used as an handler anywhere go expects a handler
```go
type myHandler struct {}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
```
myHandler statisfies the handler and can be used where a handler is expected.

However we have another important type in Go
```go
type HandlerFunc func (ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```
## Understanding ServeMux
We register a route like this; `http.HandleFunc("/", getRoot)`, where `"/"` is the path and `getRoot` is the handler. 

Doing it in the way shown above, registers the route in the default `ServeMux` and we almost never use the default in production.

A `ServeMux` is simply a router, it is what tells request where to go, for it to be handled.

The ServeMux is able to do its routing using a mechanism known as a multiplexer. A systematic process of taking in mulitple requests from clients and directing one single input to a single output the registered route to be handled.

Without the mutliplexer you would have to write multiple if else block matching the URL which isn't scalable.

You should know by now that if the multiplexer would do the routing for us, it needs to know what all our paths/patterns are, hence we need to store the somewhere after registering them. The place it gets stored is called the `DefaultServeMux`. A global router created by the std lib.

```go
// your 
http.HandleFunc("/", getRoot)
// becomes
http.DefaultServeMux.HandleFunc("/", getRoot)

// How?

// because when you write,
http.ListenAndServe(":3000", nil)
// go checks 
if handler == nil {
    handler = http.DefaultServeMux
}
// hence 
http.ListenAndServer(":3000", http.DefaultServeMux) // which means use the default global router.
```

Since we never use the default in production, lets create ours.
```go
mux := http.NewServeMux()

// registering routes on our new mux
// would be written as so
mux.HandleFunc("/", home)

// Finally we listen on our mux
http.ListenAndServe(":3000", mux)
```
Why do we never use the DefaultServeMux, it is because we do not need the same configurations for different parts of our application. Maybe we have api and cmd and development, these are different parts of the application the main func in each can scaffold it own ServeMux, while the different registered routes receive it from the main func insteading of leavning that ability for components of the app such as the Handlers to have access to them.

Since we the mutliplexer does the work for us, we dont need to use keep track of our patterns and even when the pattern doesn't match any registered pattern the multiplexer of our ServeMux takes care of the message the client sees Not Found 404.

You also spin off two server each with thier own routers or even spin off multiple server serving each serving one router
```go
// spinning of multiple servers for a single router
func main () {
    apiRoute := http.NewServeMux()

    // Register some routes
    
    go http.ListenAndServe(":3000", apiRoute)
    http.ListenAndServe(":3001", apiRoute)
}

// spinning of mulitple servers for multiple routers
func main() {
    apiRoute := http.NewServeMux()
    adminRoute := http.NewServeMux()

    // Register some routes

    go http.ListenAndServe(":3000", adminRoute)
    http.ListenAndServe(":3001", apiRoute)
}
```
### Production pattern
```go
func main() {
    // Creating a router
    mux := http.NewServeMux()

    // Registering routes
    mux.HandleFunc("/", getRoot)
    mux.HandleFunc("/health", checkHealth)

    // Creating our server
    server := &http.Server{
        Addr: ":8080"
        handler: mux
    }

    server.ListenAndServe()
}
```
## Middleware: The heart of production HTTP servers.
Engineers see middlewares as ***pipelines through which every requests flows***.

**Incoming req => Logging => Authentication => Authorization => Rate Limiting => Handler => Response => Compression => Metrics => Client**

The pipelines all request go through. You start to think of middleware when an exact implementation in that pipelines applies to all req, instead of repeting yourself each time. 

Middleware just helps you maintain the DRY principle.

So what is middleware - it is a function that sits between your middleware and your handler.

### Middleware signature
```go
func Middleware(next http.Handler) http.Handler 
// Give me the next handler in the chain, and I'll return a new handler that wraps it
```
### Building our first Middleware
```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Incoming:", r.Method(), r.URL.path)

        next.ServerHTTP(w, r)

        log.Println("Request finished")
    })
}
// the middleware can run before and after the handler
```

### Breaking the middleware pipeline
You can break out from the pipeline if the request by returning early from a middleware before exceuting the wrapped handler.
```go
func Authentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authrization")

        if token == "" {
            return
        }

        next.ServeHTTP(w,r)
    })
}
```

### Chaining Middleware
```go
middlewares := RequestMiddleware(LoggingMiddleware(AuthenticationMiddle()))
// or 
handler := RequestMiddleware()
handler = LoggingMiddleware(handler)
handler = AuthenticationMiddleware(hanlder)
```

### Other important notes
1. Ordering matters in middleware orchestration.


## Understanding `http.Server`
```go
// this is not the server, it is simply a helper
http.ListenAndServe(":3000", nil)

// this is the server
server := &http.Server{
    Addr: ":3000"
    handler: mux
}
server.ListenAndServe()

// the real server object is [http.Server]
type Server struct {
    // Addr: always load address from a config file never hardcode it. 
    // Tells the server which network address we should listen on.
    // we would usually use 0.0.0.0:port meaning listen on all ip this machine has.
    Addr string

    //this is the entry point into your application everything flow through here
    Handler Handler

    // set timeout for reading request body, prevents slowris attack
    ReadTimeout time.Duration

    // set timeout for reading request header, prevents sloworis attack
    ReadHeaderTimeout time.Duration

    // set timeout for client to read response, prevents slowris attack
    WriteTimeout time.Duration

    // modern browsers resuse connections, this process is called keep-alive
    // but we shouldn't keep it alive for open for long, reason is so we can transfer 
    // resources for other customers.
    IdleTimeout time.Duration

    // prevents the server from receiving large Headers from clients
    MaxHeaderBytes int

    ErrorLog *log.Logger

    TLSConfig *tls.Config
}


```
### Production sample
```go
server := &http.Server{
    Addr: ":" + cfg.Port
    Handler: mux

    ReadTimeout: 10 * time.Second,
    ReadHeaderTimeout: 5 * time.Second,
    IdleTimeout: 120 * time.Second,
    WriteTimeout: 30 * time.Second,

    MaxHeaderByte: 1 << 20,
}
```

Mind you these values is not a one size fits it all, based on the service you are building, you should be able to play around with the values for what makes sense for your application.


## Graceful Shutdown
Will your server stop gracefully in the middle of someone's request.

### What happens when you stop your server?
**No**. The Go runtime didn't stop your app. Your OS sent your process a signal. 

On Linux or MacOS. pressing Ctrl + C sends `SIGINT`. If your app doesn't handle that signal the OS terminates it.

### What are signal?
They are like notifications sent to a process by your OS. 

Some common signals

1. `SIGINT`: User pressed ctrl + c
2. `SIGTERM`: terminate gracefully. this is the one production apps respond to.
3. `SIGKILL`: stop immediately cannot be intercepted.
4. `SIGHUP`: configuration changed.

### Closing the server
If you care about graceful shutdown. Then you should avoid `server.Close()`. This closes all connection and stops the server immediately. So what should we use? `server.Shutdown()` it stops accepting requests and wait for the running processes to finish, closes all connections and stops the server.
```go
// create 5 sec timeout
ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
// releases resources when no longer needed, instead of waiting for timer to elapse
defer cancel()

server.Shutdown(ctx)
// the context tells us how long we should wait
```

### Production pattern
```go
func main() {
    // creating a router
    mux := http.NewServeMux()

    // registering a route
    mux.HandleFunc("/", getRoot)

    // creating a server
    server := &http.Server{
        Addr: ":8080"
        Handler: mux
        ReadTimeout: 10*time.Second
        ReadHeaderTimeout: 5*time.Second
        WriteTimeout: 10*time.Second
        IdleTimeout: 120*time.Second
        MaxHeaderByte: 1 >> 20
    }

    // Listen for requests
    go func(){
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Println(err)
            os.Exit(1)
        }
    }()

    // creating a channel
    quit := make(chan os.Signal, 1)

    // send signal to the channel
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // Wait for termination signal
    <- quit

    ctx, cancel := context.WithTimeout(context.Background, 5*time.Second)
    defer.Close()

    if err := server.Shutdown(ctx); err != nil {
        log.Println(err)
    }
}
```

## Mastering `context.Context`
Yes one of the popular usecases for context is for cancellation. However, that is not the whole picture. With context we can tie together HTTP servers, databases, gRPC, external API's, distributed tracin, cancellation, deadlines and observability.
A context is a request lifetime.

Inside your handlers the http by default passes a context for every incoming request
```go
func handler(w http.ResponseWriter, *r http.Request) {
    ctx := r.Context()

    // you can then pass this `ctx` to
    // repository
    // service 
    // etc

    // the cancellation can then propagate down the chain
}
```
You should check for the `ctx.Err()` to log the reason for the context cancellation. Mostly `context.Cancelled` `context.DeadlineExceeded` are the expected error.

Usally you would want to cancel with `ctx.WithTimeout()` but there you can also find use cases for `ctx.WithDeadline()` and `ctx.WithCancel()`

`ctx.WithValue()` this is used to hold values to pass down the chain, but you should not overdo it. Here are some values you can store using context.

1. Request ID
2. Authenticated user identity
3. trace ID
4. deadline info
5. correlation ID

**A simple rule of thumb is, if the info is realted to this request and should end only when this request ends then it should stored in context**

if you want to pass context down the chain it should be passed as the first parameter, its just a convention.

You can look [here](../context/NOTE.md) to learn more about context

## Production Error Handling
Not all kind of error should be sent as error response to the client. Some errors should be logged and some should be sent to the client.

In Go errors are values. Every layer of your application should report its own error hence no two layer return the same kind of error. Repository return a different error than the service, than the handler and so on.

You service should never leak internal state error. Like the exact error your database returns etc. However you should log internal state error.

The lower layer returns the error while the top level request handling logs it. 
```go
ErrUserNotFound := errors.New("user not found")

if errors.Is(err, ErrUserNotFound) {}

//instead of 
if err.Error() == "user not found" {}
```

### Error response
```go
{
    "error": {
        "code": "ERROR_CODE",
        "message": "error message",
    }
}

// the client can then use the code as error message might change.
```

### Central error handler
```go
// APIError represents a structured error response
type APIError struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}

// writeJSON sends a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        slog.Error("failed to encode response", "error", err)
    }
}

// writeError sends a structured error response
func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, APIError{
        Error:   http.StatusText(status),
        Message: message,
        Code:    status,
    })
}

// usage
func handler(w http.ResponseWriter, r *http.Request) {
    // perform some handler actions

    if err != nil {
        writeError(w, 500, "An unexpected error occured")
    }
}
```

### Panic is not an error
Panic should only be used when the program has reached a point where it cannot continue. Like missing database uri. Because panic can occur it is very important that we have a recovery middleware.

## Layers in backend projects

### HTTP Layer
Things that can happen in this layer
1. Parsing request
2. validating request
3. Calling services
4. returning HTTP response

### Service Layer
Responsible for business rules

### Repository Layer
Responsible for database query (postgres, redis, etc)

___
These layers move in one direction. And the way to pass each layer is through ***Dependency Injection***. With dependency injection it just saying I receive what I need, I do not create them.
```go
// service receives repo, because it uses it internally 
// the implementation in repo is not a matter of concern to the service.
service := NewUserService(repo)
```
___
The composition of all these layers is done in the main, this is where all the compositions of layers are wired together.

Instead of passing this dependencies(layers) around separately some project pass them together as struct 
```go
type Application struct {
    Logger Logger
    DB *sql.DB
    Config Config
    Server *http.Server
}
```
You need to understand the difference between application-lifetime and request-lifetime. 

Application-lifetime objects;
1. configuration
2. Db connection
3. Logger
4. Redis client
5. message broker client
6. http server

Request-lifetime objects;
1. context.Context
2. Request body
3. Authenticated user
4. Request Id
5. trace id
6. validation results.

## Steps to crafting a production level API like a Senior
### Step 1 
think about the feature you want to build and ask all relevant business questions

### Step 2
Design the folder structure. the folder structure below is a good start.
```
user-api/

├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│
│   ├── config/
│   │      config.go
│   │
│   ├── database/
│   │      postgres.go
│   │
│   ├── models/
│   │      user.go
│   │
│   ├── repository/
│   │      user_repository.go
│   │
│   ├── service/
│   │      user_service.go
│   │
│   ├── handler/
│   │      user_handler.go
│   │
│   ├── middleware/
│   │      logging.go
│   │      recovery.go
│   │      auth.go
│   │
│   ├── router/
│   │      routes.go
│   │
│   └── response/
│          json.go
│
├── migrations/
│
├── configs/
│
├── go.mod
└── go.sum
```
### Step 3 
design the domain
```go
type User struct {
    ID string
    Name string
    ...
}
```
### Step 4
design the repository interface
all the behaviour of the designed domain

### step 5
design the service layer
contains all the business logic

### step 6
design the handler

## How I create projects now
commit `feat: bootstrap production server`
1. go mod init
2. Setup config in internal/config
3. setup the logger in internal/logger
4. setup the health check handler  in internal/handler
5. set the server in internal/server
6. set the application objct in internal/app/
7. setup the composition root main.go

commit `feat: add production PostgreSQL infrastructure`
1. setup the db config in the internal/config file.
2. setup the db in internal/database/postres.go|mongodb.go
3. update the application struct
4. Start the database in the configuration root main.go
5. Close the database after the http.Server has been shutdown

commit `feat(user): introduct user domain and business layer`
1. create a domain repo for each part of the domain. the user domain for example has 
    1. entity.go which contains the struct of business data
    2. all the error that can be returned in that domain
    3. define the operation of the domain in the domain repository
    4. write the business logic in the service
    5. Perform all storage and retrival only in the repository no busienss logics no checks

commit `refactor: introduce application bootstrap`

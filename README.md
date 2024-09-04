# Concurrency in GO

## Goroutine
  - Independently executing functions
  - Lightweight threads
  - Run on top of threads
  - Runtime optimally schedules goroutines
    ```
    go funcName()
    ```

## `sync` Package
  - `sync` package provides locks and synchronization premitives
  - All of the types in the sync package should be passed by pointer to functions

### `sync.WaitGroup`
  - Used to wait for goroutines to finish
  - Under the hood, it uses a very simple counter and an inner lock
  - The zero value of the `WaitGroup` is ready to be used
  - We initialize it with `var wg sync.WaitGroup`

**Methods**
  1. `func (wg *WaitGroup) Add(delta int)`
      - Adds given number to the inner counter
      - delta is the no. of goroutines you wish to wait for
      - Adds panic if counter becomes negative
  2. `func (wg *WaitGroup) Done()`
      - Decrements inner counter and should be used when goroutine finishes its work
  3. `func (wg *WaitGroup) Wait()`
     - Blocks goroutine from which it is invoked until the counter reaches 0

**Race Detector**
  - Detects race conditions when they occur
  - Prints out stack traces and conflicting accesses when a race is found
  - Add `-race` flag to any go command to use it
    - E.g. `go run -race main.go`

### `sync.Map`
  - Safe for concurrent use by multiple goroutines
  - Equivalent to a safe `map[interface{}]interface{}`
  - The zero value is empty and ready for use
  - Incurs performance overhead and should only be used as necessary
  
**Methods**
  - `func (m *Map) Load(key interface{}) (value interface{}, ok bool)`
    - Reads existing item from map, returns nil and false ok value if none is found
  - `func (m *Map) Store(key, value interface{}) bool`
    - Inserts or updates a new key value pair
  - `func (m *Map) Range(f func(key, value interface{}) bool)`
    - Takes a function and calls it sequentially for all the values in the map

### `sync.Mutex`
  - Brings order using locks
  - The mutex is initialized **unlocked** using `var m sync.Mutex`
  - Used around critical section to execute it atomically
  
**Methods**
  - `func (m *Mutex) Lock()`
    - Locks the mutex and will block until mutex is in unlock state
  - `func (m *Mutex) Unlock()`
    - Unlocks the mutex and allows it to be used by another goroutine

### `sync.singleflight` Package
  - Provides a duplicate function call suppression mechanism
  - Don't repeat yourself - at the same time
  - Using `singleflight.Group`, we can reduce duplicate processes. Say for example, multiple web requests generating same responses. Details from the request will be used to create `key` for the function. For multiple requests with same `key`, we can avoid creating new process which will be returning same response
  - [Medium article](https://levelup.gitconnected.com/optimize-your-go-code-using-singleflight-3f11a808324)
    ```
    var sfGroup singleflight.Group
    _, err, _ := sfGroup.Do(key, func() (interface{}, error) {
      // Implementation
    })
    ```
**Methods**
  - `func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)`
    - Do executes and returns the results of the given function, making sure that only one execution is in-flight for a given key at a time
  - `func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result`
    - DoChan is like Do but returns a channel that will receive the results when they are ready
      ```
      type Result struct {
        Val    interface{}
        Err    error
        Shared bool
      }
      ```
  - `func (g *Group) Forget(key string)`
    - Forget tells the singleflight to forget about a key. Future calls to Do for this key will call the function rather than waiting for an earlier call to complete

## Channels
  - One way to share results between goroutines is to create variables in main memory
  - Channels allow goroutines to communicate among themselves and share results when they are ready
  - No need to pass values to the shared context of the main function. The channel acts as a pass-through
  - Gorountines send values to the channel and some other goroutines receive them
  - The reserved keyword `chan` denotes a channel
  - Channel operator: `<-`
  - `->, >-, -<` are invalid
  - Channels are associated with a data type and only the declared data type can be transported on them
  - The zero value of channels `var ch chan T` is nil
  - Syntax to declare a channel of type T is `ch := make(chan T)`
  - Sending is done with `ch <- data`; the arrow points into the channel as the data travels into it
  - Receiving is done with `data := <-ch` ; the arrow points away from the channel as the data travels out of it
  - Both send and receive ops are blocking

### Unbuffered Channels
  - Zero capacity channels which require both sender and receiver to be present to successfully complete ops
  - Info is exchanged synchronously

### Buffered Channels
  - Predefined capacity channels which have the ability to store values for later processing
  - Info is exchanged asynchronously

### Channel Directions
  - Bidirectional channel: `chan T`
  - Send only channel: `chan<- T`
  - Receive only channel: `<-chan T`
  - Allowed ops are enforced by the compiler
  - Bidirectional channels are implicitly cast to unidirectional channels

### Closing Channels
  - Closing a channel signals no more values will be sent to it
  - We close channel `ch` using `close(ch)`
  - We can close only bidirectional or send only channel

### Select Statement
  - Used to wait on multiple channel ops
  - Blocks until one of the channels is ready
  - When multiple ops are ready, it selects one randomly

## Concurrency Patterns
### Signalling Work Has Been Done
  - Close an additional channel called `signal (done)` channel
  - The purpose of this channel is not to transfer info but to signal work has completed
  - Its datatype is the empty struct to take up as little memory as possible
    ```
        func doWork(input <-chan string, done <-chan struct{}) {
            for {
                select {
                    case in := <-input:
                        fmt.Println("Got some input:", in)
                    case <-done:
                        return
                }
            }
        }
    ```
**Closing Channels Only Once**
  - Attempting to close an already closed channel panics
  - While the done channel stops panics on sends on the input channel, we need to ensure the signal channel is only closed once
  - The `sync` package provides `sync.Once` to help us out with this
    ```
      func sayHelloOnce() {
          var once sync.Once
          for i:= 0; i < 10; i++ {
              once.Do(func() {
                  fmt.Println("Hello, world!")
              })
          }
      }
    ```
### Worker Pools
  - A predetermined amount of workers start up
  - All workers listen for input on a shared channel
  - The shared channel is buffered
  - The same set of workers pick up multiple pieces of work

### Contexts & Cancellations
**Context**
  - A `context.Context` is generated by `net/http` for each request
  - It is available using `ctx := req.Context()` method
  - Contexts are immutable
    - We can create new context from existing context. The old context will be the parent of the new derived context

**Cancellation**
  - Allows the system to stop doing unnecessary work
  - The context exposes three ways that a request can be cancelled
    - `context.WithCancel`
    - `context.WithDeadline`
      - Specify a time after which the context will be automatically cancelled. All derived contexts are also cancelled
    - `context.WithTimeout`
      - Similar to with deadline
    - Listen for cancellation on `<-ctx.Done()`
    ```
      func doWork(ctx context.Context, input <-chan string) {
          for {
              select {
                  case in := <-input:
                    fmt.Println("Got some input:", in)
                  case <-ctx.Done():
                    fmt.Println("Out of time!", ctx.Err())
                    return
              }
          }
      }
    ```

**Why Use Context?**
  - Pass request IDs from handlers further into the app
  - Stop expensive ops from running unnecessarily
  - Keep sys latency down using a hard stop
```
  // http context usage
  func Stats(w http.ResponseWriter, r *http.Request) {
    reqCtx := r.Context()

    // New context with timeout of 100 millis
    ctx, cancel := context.WithTimeout(reqCtx, 100*time.Millisecond)
    defer cancel()

    stats, err := repo.GetOrderStats(ctx)
    if err != nil {
      writeResponse(w, http.StatusInternalServerError, nil, err)
      return
    }

    writeResponse(w, http.StatusOK, stats, nil)
  }

  func (r repo) GetOrderStats(ctx context.Context) (models.Statistics, error) {
    select {
      case s := <-r.stats.GetStats(ctx):
        return s, nil
      case <-ctx.Done():
        return models.Statistics{}, ctx.Err()
    }
  }

  func (s *statsService) GetStats(ctx context.Context) <-chan models.Statistics {
    ...
    
    select {
      ...
      case <-ctx.Done():
        fmt.Println("Context deadline exceeded)
        return
    }
    
    ...
  }
```

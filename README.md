An application can panic for numerous reasons. Whether you have an unchecked pointer dereference or are using a dependency that panics, it is important to properly "catch" an arbitrary panic so that your appliaction can clearly shutdown.

For synchronous applications, this is trival using the built in  `recover` function. However, for applications with numerous processes running concurrently, it is slighly more complex.

Regardless of how complex your application is, the premise is simple. If a panic occurrs in any of your application's execution threads, notify all other threads of execution that they must iniatiate a clean shutdown process as the application is requesting to be shutdown.

As a prerequisite, we must contruct each concurrent thread of execution so that it is preemptable. A common way to do this, as described in detail by *Peter Bourgon* in his talk, [Go + Microservices = Go Kit](), is to use the concept of run groups.

Described simply, a run group is a wrapper for any task such that the task can be started and stopped. In this example, we'll use the following interface.

```go
type Runnable interface {
	Run() error
	Close(error) error
}
```

This interface is generic enough such that any task (i.e. HTTP server, tracer, Kafka consumer, etc) can be wrapped as a `Runnable`. The contract for this interface requires that an invocation of `Close()` **must** cause the `Run()` method to return. A simple example in an instantiation that will print a string every second.

```go
type printer struct {
	str    string
	ticker *time.Ticker
	done   chan bool
}

func (p printer) Run() error {
	for {
		select {
		case <-p.done:
			return nil
		case <-p.ticker.C:
			fmt.Println(p.str)
		}
	}
}

func (p printer) Close(err error) error {
    p.ticker.Stop()
    p.done <- true
    return nil
}
```

With each of our application's concurrent operations wrapped in this manner, we can easily catch and handle all panics. This is implemented by adding the following code to the beginning of every `Run()` method.

```go
defer func() {
    if r := recover(); r != nil {
        log.Printf("recovered from panic: %v\n", r)
    }
}()
```

The driver for our runnables, which is powered by `github.com/oklog/run` in this example, can be easily understood by the following pseudocode. 

> For each runnable, invoke `Run()` in its own go routine. If the `Run()` method of any runnable returns, for each runnable, invoke `Close()`.

Our example application will execute three runnables. A printer with the aformetioned implementation, a panicker that will panic after a specified duration **but will recover because it has the aformentioned code block**, and one that will handle os signals so that we perform a clean shutdown if a signal is received (i.e. `SIGTERM`).

```go
var g Group

g.Add(Always, NewMonitor(context.Background()))
g.Add(Always, NewPrinter("Chuck Bass"))
g.Add(os.Getenv("PANICKER") == "YES", NewPanicker(5*time.Second))

if err := g.Run(); err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
    os.Exit(1)
}
os.Exit(0)
```

Here is a sample of our example application's output.

```
$ PANICKER=YES go run ./...
2020/08/17 20:59:27 starting panicker
2020/08/17 20:59:27 starting parent process monitor
2020/08/17 20:59:27 starting Chuck Bass printer
Chuck Bass
Chuck Bass
Chuck Bass
Chuck Bass
Chuck Bass
2020/08/17 20:59:32 recovered from panic: the panicker has panicked
2020/08/17 20:59:32 closing parent process monitor
2020/08/17 20:59:32 closing Chuck Bass printer
2020/08/17 20:59:32 closing panicker
```
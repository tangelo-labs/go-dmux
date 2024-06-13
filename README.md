# go-dmux

`go-dmux` is a Go package that provides an abstraction for distributed mutexes. It allows you to create and manage distributed locks using different backends. Currently, it supports Redis, in-memory and MySQL as a backend.

## Installation

To install `go-dmux`, use the `go get` command:

```bash
go get github.com/tangelo-labs/go-dmux
```

## Usage

First, you need to create a factory for the mutexes. This factory will be responsible for creating new mutexes. Here is an example of how to create a factory that uses Redis as a backend:

```go
cfg := dmux.RedisConfig{DSN: "redis://localhost:6379/0"}
factory, err := dmux.NewRedisFactory(cfg)
if err != nil {
    log.Fatalf("Failed to create factory: %v", err)
}
```

Once you have a factory, you can create a new mutex:

```go
mu, err := factory.NewMutex(context.Background(), "my-mutex")
if err != nil {
    log.Fatalf("Failed to create mutex: %v", err)
}
```

To lock and unlock the mutex, use the `Lock` and `Unlock` methods:

```go
err = mu.Lock(context.Background())
if err != nil {
    log.Fatalf("Failed to lock: %v", err)
}

// Do some work...

err = mu.Unlock(context.Background())
if err != nil {
    log.Fatalf("Failed to unlock: %v", err)
}
```

## Testing

The package includes a test suite that you can run with the `go test` command:

```bash
go test ./...
```

## Contributing

Contributions to `go-dmux` are welcome. Please submit a pull request or create an issue to discuss the changes you want to make.

## License

`go-dmux` is licensed under the MIT License. See the `LICENSE` file for more information.

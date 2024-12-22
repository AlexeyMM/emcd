# log

Logger (wrapper over zerolog) and Interceptor, to support traceID and userID in logs.

## Installation

```bash
go get -u code.emcdtech.com/emcd/sdk/log
```

## Getting started

### Simple Logging Example

For simple logging, import the global logger package **code.emcdtech.com/emcd/sdk/log**

```go
package main

import (
    "code.emcdtech.com/emcd/sdk/log"
)

func main() {
    ctx := context.Background()
    err := log.Init(ctx) //not mandatory, call it once in config initialization
    if err!=nil{
        panic()
    }

    log.Info(ctx, "hello world")
    log.Info(ctx, "hello %s", "world")
}
```

> Note: By default log writes to `os.Stderr`
> 
> Note: By default log level picked up from LOG_LEVEL env
> 
> Note: By default time format is 2006-01-02 15:04:05.000 and output is json

#### Logging Methods

```go
type Logger interface {
	Err(ctx context.Context, err error)
	Debug(ctx context.Context, format string, v ...interface{})
	Info(ctx context.Context, format string, v ...interface{})
	Warn(ctx context.Context, format string, v ...interface{})
	Error(ctx context.Context, format string, v ...interface{})
	Fatal(ctx context.Context, format string, v ...interface{})
}
```

### Interceptor

Interceptor allows you to inject `userID` in `context.Context` and support apm tracing

### Simple Interceptor Example

Inject `grpc.UnaryInterceptor(log.ServerUnaryInterceptor)` and `apmgrpc.NewUnaryClientInterceptor()` to `grcp.NewServer()`

```go
package main

import (
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
	
    "code.emcdtech.com/emcd/sdk/log"
	"go.elastic.co/apm/module/apmgrpc/v2"
)

func main() {
	transport, err := transport.NewHTTPTransport(transport.HTTPTransportOptions{})
	if err != nil {
		log.Fatal(context.Background(), "create transport for apm: %w", err)
	}
	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:        "service name",
		ServiceVersion:     "version",
		ServiceEnvironment: "development",
		Transport:          transport,
	})
	if err != nil {
		log.Fatal(context.Background(), "create tracer: %w", err)
	}
	
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(log.ServerUnaryInterceptor),
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery(), apmgrpc.WithTracer(tracer)),
	)
}
```

For client GRPC calls, inject `grpc.WithChainUnaryInterceptor(log.ClientUnaryInterceptor)` and `apmgrpc.NewUnaryClientInterceptor()` to `grpc.DialOption`

```go
package main

import (
    "code.emcdtech.com/emcd/sdk/log"
	"go.elastic.co/apm/module/apmgrpc/v2"
)

func NewAutoPay(url string) proto.AdminServiceClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		apmgrpc.NewUnaryClientInterceptor(),
		grpc.WithChainUnaryInterceptor(log.ClientUnaryInterceptor),
	}

	connServ, err := grpc.Dial(url, opts...)
	if err != nil {
		log.Fatal(context.Background(), "fail to dial: %v", err)
	}

	client := proto.NewAdminServiceClient(connServ)

	return client
}
```
### List of fields
```
"level" - importance of log record trace|debug|info|warn|error|fatal
"time" - date/time of log record 2006-01-02 15:04:05.000
"message" - text information about event
"user.id" - user id
"trace.id" - uniq identifier for HTTP request. Usually created on HTTP server and transfer to another service through context. If trace ID is empty - create new.
"transaction.id" - id of transaction
"service" - source of service
```

### Log requirements
https://docs.google.com/document/d/1yqGXd9H36Rwjza7wZRHFOd2CIv2ncwH5WoEC1tH6HFE/edit
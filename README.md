# grpctest

A small testing library for [gRPC-Go](https://github.com/grpc/grpc-go) servers

## Description

grpctest provides a test server and client connection.
The server response from a given client request is obtained with trivial setup, reducing the overhead of testing a gRPC service.

## Usage

This test server is intended to be used along side the protocol compiler plugins,
[protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) and
[protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) to generate the relevant code.
Familiarity with this process is assumed in the instructions below to further simplify setup.
The following uses the [gRPC quickstart](https://grpc.io/docs/languages/go/quickstart/) greeter server as an example.

### Set up the test

> [example code](https://github.com/grpc/grpc-go/tree/master/examples/helloworld)

Install the necessary prerequistes and create a test file. Import grpctest.

```go
import "github.com/vision-cli/grpctest"
```

Define a test function. A test server can be initialised with context.

```go
func TestSayHello(t *testing.T) {
	ctx := context.Background()
	s := grpctest.NewServer().WithContext(ctx)
	defer s.Close()

}
```

### Initialise the client

Import the relevant generated protobuf package.

```go
import pb "google.golang.org/grpc/examples/helloworld/helloworld"
```

Register the greeter server using the generated function with `s.RunServer()`.

```go
s.RunServer(t, func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &server{})
	})
```

Initialise a new Greeter client with `s.ClientConn()`

```go
client := pb.NewGreeterClient(s.ClientConn(t))
```

### Test a request

With all the necesssary setup complete, the client can now make a request.

```go
reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})
```

The reply and any errors can be tested accordingly.

```go
if err != nil {
	t.Fatalf("SayHello failed: %v", err)
}
expected := "Hello Alice"
if actual := reply.GetMessage(); actual != expected {
	t.Errorf(`got "%s", want "%s"`, actual, expected)
}
```

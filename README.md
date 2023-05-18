# Activity Interceptor for gRPC

The Activity Interceptor is a tool designed to enhance the observability and debugging capabilities of gRPC services. It provides a mechanism to extract and process valuable information from incoming gRPC requests and their corresponding responses. This information includes the caller's IP address, the time of the request and response, the operation ID, the RPC path, and the status code of the response.

The interceptor is implemented as a middleware that wraps around your gRPC service handlers. This allows it to intercept incoming requests, extract the necessary information, and then pass the request along to the actual service handler. After the handler has processed the request and generated a response, the interceptor can then extract additional information from the response before it is sent back to the client.

## Why is it Useful?

The Activity Interceptor is a valuable tool for several reasons:

1. **Observability**: It provides a way to observe and log the activity of your gRPC services in real-time. This can be useful for monitoring the health and performance of your services, as well as for identifying and diagnosing issues.

2. **Debugging**: By logging detailed information about each request and response, the interceptor can provide valuable insights that can aid in debugging. For example, if a request is failing, the interceptor can provide information about the request (such as the caller's IP and the RPC path) that can help identify the source of the problem.

3. **Customizability**: The interceptor is highly customizable. You can specify a custom callback function to be called for each request and response, allowing you to process the activity information in any way you see fit. You can also customize the operation ID header and the function used to generate operation IDs.

## Usage

To use the Activity Interceptor, you first need to create an instance of it by calling the `NewActivityInterceptor` function and passing in an `ActivityInterceptorOpts` struct. This struct allows you to specify a custom callback function, operation ID header, and operation ID generation function. If you don't provide these options, the interceptor will use sensible defaults.

Once you have an instance of the interceptor, you can use its `UnaryServerInterceptor` method to get a gRPC unary server interceptor. This interceptor can be added to your gRPC server using the `grpc.UnaryInterceptor` server option.

Here's an example of how to use the Activity Interceptor:

```go
opts := interceptor.ActivityInterceptorOpts{
	Callback: func(ctx context.Context, info *interceptor.ActivityInfo) {
		// Process the activity info here...
	},
  // Optional.
	OperationIDHeader: "x-custom-operation-id",
  // Optional.
	OperationIDFunc: func() string {
		// Generate a custom operation ID here...
	},
}
interceptor := interceptor.NewActivityInterceptor(opts)
grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()))
```

## Models

The Activity Interceptor exports two main models:

1. `ActivityInfo`: This struct contains the information extracted from a request and its corresponding response. It includes the caller's IP address, the time of the request and response, the operation ID, the RPC path, and the status code of the response.

2. `ActivityInterceptorOpts`: This struct allows you to customize the behavior of the Activity Interceptor. You can specify a custom callback function, operation ID header, and operation ID generation function.

## Callback Function

The callback function is a key part of the Activity Interceptor. It is called for each request and response, and is passed a context and an `ActivityInfo` struct containing the information extracted from the request and response.

The context passed to the callback function is the same context that is passed to the gRPC service handler. This means it contains any metadata associated
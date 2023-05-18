package interceptor_test

import (
	"context"
	"net"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	interceptor "github.com/problem-company-toolkit/activity-interceptor"
)

var _ = Describe("Interceptor", func() {
	var (
		ctx     context.Context
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	)

	BeforeEach(func() {
		ctx = context.Background()
		info = &grpc.UnaryServerInfo{FullMethod: "/test.v1/TestMethod"}
		handler = func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, status.Error(codes.NotFound, "not found")
		}
	})

	Context("when the request is processed", func() {
		It("should extract the correct RPCPath", func() {
			var extractedRPCPath string
			callback := func(ctx context.Context, info *interceptor.ActivityInfo) {
				extractedRPCPath = info.RPCPath
			}
			interceptor := interceptor.NewActivityInterceptor(interceptor.ActivityInterceptorOpts{Callback: callback})

			interceptor.UnaryServerInterceptor()(ctx, nil, info, handler)

			Expect(extractedRPCPath).To(Equal("/test.v1/TestMethod"))
		})

		It("should extract the correct IP", func() {
			var extractedIP string
			callback := func(ctx context.Context, info *interceptor.ActivityInfo) {
				extractedIP = info.CallerIP
			}
			interceptor := interceptor.NewActivityInterceptor(interceptor.ActivityInterceptorOpts{Callback: callback})

			ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1")}})
			interceptor.UnaryServerInterceptor()(ctx, nil, info, handler)

			Expect(extractedIP).To(Equal("127.0.0.1"))
		})

		It("should extract the correct status code", func() {
			var extractedStatusCode int
			callback := func(ctx context.Context, info *interceptor.ActivityInfo) {
				extractedStatusCode = info.StatusCode
			}
			interceptor := interceptor.NewActivityInterceptor(interceptor.ActivityInterceptorOpts{Callback: callback})

			interceptor.UnaryServerInterceptor()(ctx, nil, info, handler)

			Expect(extractedStatusCode).To(BeEquivalentTo(codes.NotFound))
		})

		It("should not generate a new operation ID if the operation ID header already has a value", func() {
			var extractedOperationID string
			callback := func(ctx context.Context, info *interceptor.ActivityInfo) {
				extractedOperationID = info.OperationID
			}
			interceptor := interceptor.NewActivityInterceptor(interceptor.ActivityInterceptorOpts{Callback: callback})

			ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("x-operation-id", "existing-operation-id"))
			interceptor.UnaryServerInterceptor()(ctx, nil, info, handler)

			Expect(extractedOperationID).To(Equal("existing-operation-id"))
		})
	})
})

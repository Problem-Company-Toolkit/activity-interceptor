package interceptor

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	_DEFAULT_OPERATION_ID_HEADER = "x-operation-id"
)

type ActivityInfo struct {
	CallerIP     string
	RequestTime  time.Time
	ResponseTime time.Time
	OperationID  string
	RPCPath      string
	StatusCode   int
}

type callbackFunc = func(context.Context, *ActivityInfo)
type operationIDFunc = func() string

type ActivityInterceptor struct {
	operationIDHeader   string
	callback            callbackFunc
	generateOperationID operationIDFunc
}

type ActivityInterceptorOpts struct {

	// Optional. Specifies an action to do when receiving the request. By default just logs to stdout.
	Callback callbackFunc

	// Optional. By default it's "x-operation-id"
	OperationIDHeader string

	// Optional. By default, we generate UUIDs.
	OperationIDFunc operationIDFunc
}

func NewActivityInterceptor(opts ActivityInterceptorOpts) *ActivityInterceptor {
	if opts.Callback == nil {
		opts.Callback = func(ctx context.Context, info *ActivityInfo) {
			log.Printf("%+v", info)
		}
	}
	if opts.OperationIDHeader == "" {
		opts.OperationIDHeader = _DEFAULT_OPERATION_ID_HEADER
	}
	if opts.OperationIDFunc == nil {
		opts.OperationIDFunc = uuid.NewString
	}

	return &ActivityInterceptor{
		operationIDHeader:   opts.OperationIDHeader,
		callback:            opts.Callback,
		generateOperationID: opts.OperationIDFunc,
	}
}

func (a *ActivityInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		activityInfo := &ActivityInfo{
			RequestTime: time.Now(),
			RPCPath:     info.FullMethod,
		}

		if p, ok := peer.FromContext(ctx); ok {
			activityInfo.CallerIP = p.Addr.(*net.TCPAddr).IP.String()
		}

		md, _ := metadata.FromIncomingContext(ctx)
		if operationID := md.Get(a.operationIDHeader); len(operationID) > 0 {
			activityInfo.OperationID = operationID[0]
		} else {
			activityInfo.OperationID = a.generateOperationID()
		}

		resp, err = handler(ctx, req)

		activityInfo.ResponseTime = time.Now()
		if err != nil {
			activityInfo.StatusCode = int(status.Code(err))
		}

		a.callback(ctx, activityInfo)

		return resp, err
	}
}

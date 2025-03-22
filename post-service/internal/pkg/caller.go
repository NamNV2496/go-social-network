package pkg

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type callOption struct {
	timeout time.Duration
}

type CallOption func(*callOption)

func GRPCCaller[request, response, input, output any](
	ctx context.Context,
	req *request,
	fn func(context.Context, *input, ...grpc.CallOption) (*output, error),
	opts []CallOption,
	grpcOpts ...grpc.CallOption,
) (*response, error) {
	opt := callOption{}
	for _, o := range opts {
		o(&opt)
	}
	var in input
	if err := Copy(&in, req); err != nil {
		return nil, err
	}
	callCtx := ctx
	if opt.timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, opt.timeout)
		defer cancel()
		callCtx = timeoutCtx
	}
	out, err := fn(callCtx, &in, grpcOpts...)
	if err != nil {
		return nil, err
	}
	var resp response
	if err := Copy(&resp, out); err != nil {
		return nil, err
	}
	return &resp, err
}

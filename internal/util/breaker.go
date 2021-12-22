package util

import (
	"context"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func WithCircuitBreaker(cb *gobreaker.CircuitBreaker, openHandler func(), closeHandler func()) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		_, err := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

		if err == gobreaker.ErrOpenState {
			openHandler()
		} else if err == nil && cb.State() == gobreaker.StateClosed {
			closeHandler()
		}

		return err
	}
}

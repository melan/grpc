package util

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errLimitExceeded = status.Errorf(codes.ResourceExhausted, "Limit exhausted, please try later")
)

type AllowLimiter interface {
	Allow() bool
}

type WaitLimiter interface {
	Wait(ctx context.Context) error
}

func NewErrorLimiter(limiter AllowLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !limiter.Allow() {
			fmt.Println("limit exceeded")
			return nil, errLimitExceeded
		}

		return handler(ctx, req)
	}
}

func NewWaitLimiter(limiter WaitLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err := limiter.Wait(ctx); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

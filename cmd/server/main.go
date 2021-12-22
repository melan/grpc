package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/melan/grpc/internal/server"
	"github.com/melan/grpc/internal/util"
	apiv1 "github.com/melan/grpc/pkg/api/v1"
	apiv2 "github.com/melan/grpc/pkg/api/v2"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port  = flag.Int("port", 50051, "The server port")
	token = flag.String("token", "", "Security token")
	limit = flag.Int("limit", 10, "Requests limit")

	key  = flag.String("key", "key.pem", "Private key file")
	cert = flag.String("cert", "cert.pem", "Public cert file")

	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if token != nil && *token != "" && !valid(md["authorization"], *token) {
		return nil, errInvalidToken
	}

	return handler(ctx, req)
}

func valid(authorization []string, expectedToken string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == expectedToken
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		panic(fmt.Sprintf("can't start listener on the port %d: %s", *port, err))
	}

	sv1 := &server.V1PingServer{}
	sv2 := &server.V2PingServer{}

	var opts []grpc.ServerOption

	grpc.EnableTracing = true

	if token != nil && *token != "" {
		creds, err := credentials.NewServerTLSFromFile(*cert, *key)
		if err != nil {
			panic(fmt.Sprintf("can't load server key pair. key %q, cert %q: %s", *key, *cert, err))
		}
		opts = append(opts, grpc.Creds(creds))

		limiter := rate.NewLimiter(rate.Limit(*limit), *limit)
		fmt.Printf("server will be limeted to %d calls per second\n", *limit)

		opts = append(opts,
			grpc.ChainUnaryInterceptor(
				util.NewWaitLimiter(limiter),
				authInterceptor,
			),
		)
	}

	grpcS := grpc.NewServer(opts...)

	apiv1.RegisterPingServer(grpcS, sv1)
	apiv2.RegisterPingServer(grpcS, sv2)

	if err := grpcS.Serve(listener); err != nil {
		panic(fmt.Sprintf("can't start GRPC server: %s", err))
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/melan/grpc/internal/util"
	apiv2 "github.com/melan/grpc/pkg/api/v2"
	"github.com/sony/gobreaker"
	"time"

	apiv1 "github.com/melan/grpc/pkg/api/v1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	host  = flag.String("host", "localhost", "Server host")
	port  = flag.Int("port", 50051, "Server port")
	cert  = flag.String("cert", "cert.pem", "Server cert")
	name  = flag.String("name", "foo", "Name to send to the server")
	id    = flag.Int("id", 100, "id for the request")
	token = flag.String("token", "", "Security token")
	retry = flag.Int("retry", 0, "repeat the calls upto this number of times")
)

func main() {
	flag.Parse()

	var grpcOpts []grpc.DialOption

	if token != nil && *token != "" {
		creds, err := credentials.NewClientTLSFromFile(*cert, "")
		if err != nil {
			panic(fmt.Sprintf("can't load client cert %q %s", *cert, err))
		}

		oauthToken := oauth2.Token{TokenType: "Bearer", AccessToken: *token}
		grpcOpts = append(grpcOpts,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&oauthToken)),
		)
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "MainBreaker",
		Timeout: 1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 0
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("%s CB %s change state from %s to %s\n", time.Now().Format(time.StampMilli), name, from.String(), to.String())
		},
	})

	delay := 1000

	grpcOpts = append(grpcOpts,
		grpc.WithChainUnaryInterceptor(
			util.WithCircuitBreaker(cb,
				func() {
					time.Sleep(1 * time.Second)
					delay = delay * 2
					if delay > 10000 {
						delay = 10000
					}
				},
				func() {
					delay = int(float32(delay) * 0.9)
					if delay < 100 {
						delay = 100
					}
				}),
		))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), grpcOpts...)
	if err != nil {
		panic(fmt.Sprintf("can't call GRPC server: %s", err))
	}

	clientV1 := apiv1.NewPingClient(conn)
	clientV2 := apiv2.NewPingClient(conn)

	for i := 0; i < *retry; i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		if response, err := clientV1.Ping(context.TODO(), &apiv1.PingRequest{Name: *name}); err != nil {
			fmt.Printf("failed to call v1 GRPC server: %s\n", err)
		} else {
			fmt.Printf("%s Response from the v1 server: %s\n", time.Now().Format(time.RFC3339Nano), response.Phrase)
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)

		if response, err := clientV2.Ping(context.TODO(), &apiv2.PingRequest{ID: int32(*id)}); err != nil {
			fmt.Printf("failed to call v2 GRPC server: %s\n", err)
		} else {
			fmt.Printf("%s Response from the v2 server: %s\n", time.Now().Format(time.RFC3339Nano), response.Phrase)
		}

		fmt.Printf("Delay is %d\n", delay)
	}
}

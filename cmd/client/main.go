package main

import (
	"context"
	"flag"
	"fmt"
	apiv2 "github.com/melan/grpc/pkg/api/v2"

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

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), grpcOpts...)
	if err != nil {
		panic(fmt.Sprintf("can't call GRPC server: %s", err))
	}

	clientV1 := apiv1.NewPingClient(conn)
	if response, err := clientV1.Ping(context.TODO(), &apiv1.PingRequest{Name: *name}); err != nil {
		panic(fmt.Sprintf("failed to call v1 GRPC server: %s", err))
	} else {
		fmt.Printf("Response from the v1 server: %s\n", response.Phrase)
	}

	clientV2 := apiv2.NewPingClient(conn)
	if response, err := clientV2.Ping(context.TODO(), &apiv2.PingRequest{ID: int32(*id)}); err != nil {
		panic(fmt.Sprintf("failed to call v2 GRPC server: %s", err))
	} else {
		fmt.Printf("Response from the v2 server: %s\n", response.Phrase)
	}
}

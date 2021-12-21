package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/melan/grpc/pkg/api"
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

	client := api.NewPingClient(conn)
	response, err := client.Ping(context.TODO(), &api.PingRequest{Name: *name})
	if err != nil {
		panic(fmt.Sprintf("failed to call GRPC server: %s", err))
	}

	fmt.Printf("Response from the server: %s\n", response.Phrase)
}

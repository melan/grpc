package server

import (
	"context"
	"fmt"

	"github.com/melan/grpc/pkg/api"
)

var _ api.PingServer = &PingServer{}

type PingServer struct {
	api.UnimplementedPingServer
}

func (p *PingServer) Ping(_ context.Context, request *api.PingRequest) (*api.PingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("can't work with a nil request")
	}

	name := request.Name

	return &api.PingResponse{Phrase: fmt.Sprintf("Hello %s", name)}, nil
}

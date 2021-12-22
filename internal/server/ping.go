package server

import (
	"context"
	"fmt"

	apiv1 "github.com/melan/grpc/pkg/api/v1"
	apiv2 "github.com/melan/grpc/pkg/api/v2"
)

var _ apiv1.PingServer = &V1PingServer{}

type V1PingServer struct {
	apiv1.UnimplementedPingServer
}

func (p *V1PingServer) Ping(_ context.Context, request *apiv1.PingRequest) (*apiv1.PingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("can't work with a nil request")
	}

	name := request.Name

	return &apiv1.PingResponse{Phrase: fmt.Sprintf("Hello %s", name)}, nil
}

type V2PingServer struct {
	apiv2.UnimplementedPingServer
}

func (p *V2PingServer) Ping(_ context.Context, request *apiv2.PingRequest) (*apiv2.PingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("can't work with a nil request")
	}

	id := request.ID

	return &apiv2.PingResponse{Phrase: fmt.Sprintf("Hello, ID %d", id)}, nil
}

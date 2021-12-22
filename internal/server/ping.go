package server

import (
	"context"
	"fmt"

	apiv1 "github.com/melan/grpc/pkg/api/v1"
)

var _ apiv1.PingServer = &PingServer{}

type PingServer struct {
	apiv1.UnimplementedPingServer
}

func (p *PingServer) Ping(_ context.Context, request *apiv1.PingRequest) (*apiv1.PingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("can't work with a nil request")
	}

	name := request.Name

	return &apiv1.PingResponse{Phrase: fmt.Sprintf("Hello %s", name)}, nil
}

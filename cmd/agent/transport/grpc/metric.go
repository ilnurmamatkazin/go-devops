package grpc

import (
	"context"
	"fmt"

	g "github.com/ilnurmamatkazin/go-devops/cmd/server/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GRPCClient) SendMetrics(data []byte) error {
	req := &g.GRPCMetrics{
		// Metrics: string(data),
	}

	_, err := c.mc.SendMetrics(context.Background(), req)
	if err != nil {
		return fmt.Errorf("Error on GetUsers rpc call: %v \n", err)
	}

	return nil
}

func (c *GRPCClient) SendMetric(m g.MetricsClient) (*emptypb.Empty, error) {

	return nil, nil
}

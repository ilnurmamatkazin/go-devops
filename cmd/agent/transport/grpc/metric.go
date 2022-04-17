package grpc

import (
	"context"

	g "github.com/ilnurmamatkazin/go-devops/cmd/agent/grpc"

	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) SendMetrics(ctx context.Context, metrics string) error {
	md := metadata.Pairs("key", c.cfg.Key)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &g.GRPCMetric{
		Metric: metrics,
	}

	_, err := c.mc.SendMetrics(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *GRPCClient) SendMetric(ctx context.Context, metrics []string) error {
	md := metadata.Pairs("key", c.cfg.Key)
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := c.mc.SendMetric(ctx)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		stream.Send(&g.GRPCMetric{
			Metric: metric,
		})
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		return err
	}

	return nil
}

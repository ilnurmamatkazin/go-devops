package grpc

import (
	"context"

	g "github.com/ilnurmamatkazin/go-devops/cmd/agent/grpc"

	"google.golang.org/grpc/metadata"
)

// SendMetrics отправка массива метрик в одном запросе
func (c *GRPCClient) SendMetrics(ctx context.Context, metrics string) error {
	// передача значения ключа в метаданных
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

// SendMetric потоковая отправка атомарных метрик в одном сеансе
func (c *GRPCClient) SendMetric(ctx context.Context, metrics []string) error {
	// передача значения ключа в метаданных
	md := metadata.Pairs("key", c.cfg.Key)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// создание потока
	stream, err := c.mc.SendMetric(ctx)
	if err != nil {
		return err
	}

	// цикл отправки атомарных метрик
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

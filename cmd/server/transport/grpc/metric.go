package grpc

import (
	"context"
	"encoding/json"

	cr "github.com/ilnurmamatkazin/go-devops/cmd/server/crypto"
	g "github.com/ilnurmamatkazin/go-devops/cmd/server/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SendMetric(ctx context.Context, req *g.GRPCMetric) (*emptypb.Empty, error) {
	var (
		metric models.Metric
		err    error
	)

	b := []byte(req.MetricType)

	decodeBody, err := cr.Decrypt(s.Cfg.PrivateKey, b)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(decodeBody, &metric); err != nil {
		return nil, err
	}

	if err = s.Service.SetMetric(metric); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *server) SendMetrics(ctx context.Context, req *g.GRPCMetrics) (*emptypb.Empty, error) {
	var (
		metric []models.Metric
		err    error
	)

	b := []byte(req.Metrics[0].MetricType)

	decodeBody, err := cr.Decrypt(s.Cfg.PrivateKey, b)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(decodeBody, &metric); err != nil {
		return nil, err
	}

	if err = s.Service.SetArrayMetrics(metric); err != nil {
		return nil, err
	}

	return nil, nil
}

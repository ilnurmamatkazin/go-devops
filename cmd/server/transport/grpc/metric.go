package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	cr "github.com/ilnurmamatkazin/go-devops/cmd/server/crypto"
	g "github.com/ilnurmamatkazin/go-devops/cmd/server/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SendMetric(stream g.Metrics_SendMetricServer) error {
	var (
		metric models.Metric
	)

	md, ok := metadata.FromIncomingContext(stream.Context())
	fmt.Println("!!!!!!!!!!!!! md, ok", md, ok, md.Get("key")[0])

	if md.Get("protocol")[0] != s.Cfg.Key {
		return fmt.Errorf("для протокола grpc не совпадают ключи")
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&emptypb.Empty{})
			}

			return fmt.Errorf("Ошибка чтения клиентского запроса: %v", err)
		}

		// Get the title, pages and year fields from the req
		strMetric := req.GetMetric()

		b := []byte(strMetric)

		decodeBody, err := cr.Decrypt(s.Cfg.PrivateKey, b)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(decodeBody, &metric); err != nil {
			return err
		}

		if err = s.Service.SetMetric(metric); err != nil {
			return err
		}

	}
}

func (s *server) SendMetrics(ctx context.Context, req *g.GRPCMetric) (*emptypb.Empty, error) {
	var (
		metric []models.Metric
		err    error
	)

	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println("!!!!!!!!!!!!! md, ok", md, ok, md.Get("key")[0])

	if md.Get("protocol")[0] != s.Cfg.Key {
		return nil, fmt.Errorf("для протокола grpc не совпадают ключи")
	}

	b := []byte(req.Metric)

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

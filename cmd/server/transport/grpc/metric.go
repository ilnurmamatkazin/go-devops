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

// SendMetric обработка потоковой передачи атомарных метрик.
// Идея заключается в том, что бы получить все метрики в одном сеансе.
func (s *server) SendMetric(stream g.Metrics_SendMetricServer) error {
	var (
		metric models.Metric
	)

	// получаем ключ из метаданных
	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok || md.Get("key")[0] != s.Cfg.Key {
		return fmt.Errorf("для протокола grpc не совпадают ключи")
	}

	// в цикле читаем все атомарные метрики метрики
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&emptypb.Empty{})
			}

			return fmt.Errorf("ошибка чтения клиентского запроса %v", err)
		}

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

// SendMetrics стандартный механизм запрос-ответ.
// В запросе приходит массив метрик, который целиком обрабатывается.
func (s *server) SendMetrics(ctx context.Context, req *g.GRPCMetric) (*emptypb.Empty, error) {
	var (
		metrics []models.Metric
		err     error
	)

	// получаем ключ из метаданных
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || md.Get("key")[0] != s.Cfg.Key {
		return new(emptypb.Empty), fmt.Errorf("для протокола grpc не совпадают ключи")
	}

	b := []byte(req.Metric)

	decodeBody, err := cr.Decrypt(s.Cfg.PrivateKey, b)
	if err != nil {
		return new(emptypb.Empty), err
	}

	if err = json.Unmarshal(decodeBody, &metrics); err != nil {
		return new(emptypb.Empty), err
	}

	if err = s.Service.SetArrayMetrics(metrics); err != nil {
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

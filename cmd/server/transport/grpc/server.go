package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	g "github.com/ilnurmamatkazin/go-devops/cmd/server/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
)

type server struct {
	g.UnimplementedMetricsServer
	Cfg     models.Config
	Service *service.Service
}

// StartGRPC создаем стандартный gRPC сервер
func StartGRPC(cfg models.Config, service *service.Service) {
	log.Println("Starting server..")

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Printf("Unable to listen on port %s: %v", cfg.GRPCPort, err)
	}

	s := grpc.NewServer()
	g.RegisterMetricsServer(s, &server{Cfg: cfg, Service: service})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fails to serve: %v", err)
	}
}

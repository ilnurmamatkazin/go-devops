package grpc

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	g "github.com/ilnurmamatkazin/go-devops/cmd/agent/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

type GRPCClient struct {
	con *grpc.ClientConn
	cfg models.Config
	mc  g.MetricsClient
}

// Создание стандартного gRPC клиента
func NewGRPCClient(cfg models.Config) (*GRPCClient, error) {
	log.Println("GRPC Client..")

	var err error

	client := GRPCClient{cfg: cfg}

	client.con, err = grpc.Dial(cfg.AddressGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error connecting: %v \n", err)
	}

	client.mc = g.NewMetricsClient(client.con)

	return &client, nil
}

func (c *GRPCClient) Close() {
	c.con.Close()
}

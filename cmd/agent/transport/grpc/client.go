package grpc

import (
	"fmt"

	"google.golang.org/grpc"

	g "github.com/ilnurmamatkazin/go-devops/cmd/server/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

type GRPCClient struct {
	con *grpc.ClientConn
	cfg models.Config
	mc  g.MetricsClient
}

func NewGRPCClient(cfg models.Config) (*GRPCClient, error) {
	fmt.Println("GRPC Client..")

	var err error

	client := GRPCClient{cfg: cfg}
	opts := grpc.WithInsecure()

	client.con, err = grpc.Dial("localhost:3000", opts)
	if err != nil {
		return nil, fmt.Errorf("Error connecting: %v \n", err)
	}

	// defer con.Close()
	client.mc = g.NewMetricsClient(client.con)
	// getUsers(c)

	return &client, nil
}

func (c *GRPCClient) Close() {
	c.con.Close()
}

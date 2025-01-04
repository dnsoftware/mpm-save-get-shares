package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type GRPCCoinStorage struct {
	conn   *grpc.ClientConn
	client MinersServiceClient
}

func NewCoinStorage(conn *grpc.ClientConn) (*GRPCCoinStorage, error) {

	client := NewMinersServiceClient(conn)

	return &GRPCCoinStorage{
		client: client,
		conn:   conn,
	}, nil
}

func (g *GRPCCoinStorage) GetCoinIDByName(ctx context.Context, coin string) (int64, error) {

	state := g.conn.GetState()
	log.Printf("Connection state: %v", state.String())

	resp, err := g.client.GetCoinIDByName(ctx, &GetCoinIDByNameRequest{
		Coin: coin,
	})

	if err != nil {
		return 0, err
	}

	return resp.Id, err
}

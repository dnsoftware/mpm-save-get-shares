package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/dnsoftware/mpm-save-get-shares/internal/adapter/grpc/proto"
)

type GRPCCoinStorage struct {
	conn   *grpc.ClientConn
	client proto.MinersServiceClient
	//jwtProcessor JWTProcessor
	//jwtToken     string
}

func NewCoinStorage(conn *grpc.ClientConn /*, jwtProcessor JWTProcessor*/) (*GRPCCoinStorage, error) {

	client := proto.NewMinersServiceClient(conn)

	//token, err := jwtProcessor.GetActualToken()
	//if err != nil {
	//	return nil, err
	//}

	return &GRPCCoinStorage{
		client: client,
		conn:   conn,
		//jwtProcessor: jwtProcessor,
		//jwtToken:     token,
	}, nil
}

func (g *GRPCCoinStorage) GetCoinIDByName(ctx context.Context, coin string) (int64, error) {

	resp, err := g.client.GetCoinIDByName(ctx, &proto.GetCoinIDByNameRequest{
		Coin: coin,
	})

	if err != nil {
		return 0, err
	}

	return resp.Id, err
}

package client

import (
	"fmt"
	"github.com/Diode222/GomokuGameReferee/conf"
	pb "github.com/Diode222/GomokuGameReferee/probo"
	"google.golang.org/grpc"
)

func NewPlayerServiceClient(port string) (pb.MakePieceServiceClient, error) {
	sreviceAddr := fmt.Sprintf("%s:%s", conf.IP, port)
	conn, err := grpc.Dial(sreviceAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewMakePieceServiceClient(conn), nil
}

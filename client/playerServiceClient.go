package client

import (
	"fmt"
	"github.com/Diode222/GomokuGameReferee/conf"
	pb "github.com/Diode222/GomokuGameReferee/probo"
	"google.golang.org/grpc"
	"log"
)

func NewPlayerServiceClient(port string) (pb.MakePieceServiceClient, error) {
	sreviceAddr := fmt.Sprintf("%s:%s", conf.IP, port)
	conn, err := grpc.Dial(sreviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(fmt.Sprintf("Connect to %s failed", sreviceAddr))
		return nil, err
	}

	return pb.NewMakePieceServiceClient(conn), nil
}

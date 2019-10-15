package client

import (
	"github.com/Diode222/GomokuGameReferee/conf"
	pb "github.com/Diode222/GomokuGameReferee/probo"
	"github.com/Diode222/etcd_service_discovery/etcdservice"
)

func NewPlayerServiceClient(serviceName string) pb.MakePieceServiceClient {
	return etcdservice.NewServiceManager(conf.ETCD_ADDR).GetClient(serviceName, pb.NewMakePieceServiceClientWrapper).(pb.MakePieceServiceClient)
}

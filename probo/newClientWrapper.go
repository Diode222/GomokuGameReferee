package proto

import "google.golang.org/grpc"

func NewMakePieceServiceClientWrapper(cc *grpc.ClientConn) interface{} {
	return NewMakePieceServiceClient(cc)
}

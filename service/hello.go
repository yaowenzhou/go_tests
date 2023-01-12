package service

import (
	"context"
	"go_tests/pb"
)

// HelloServer TODO
type HelloServer struct {
	pb.UnimplementedHelloServer
}

// SayHello TODO
func (s *HelloServer) SayHello(ctx context.Context, in *pb.SayHelloReq) (out *pb.SayHelloResp, err error) {
	return &pb.SayHelloResp{Message: "Hello, Grpc!"}, nil
}

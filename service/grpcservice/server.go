package grpcservice

import (
	"context"
	"fmt"
	"net"

	"github.com/rlaskowski/easymotion/config"
	pb "github.com/rlaskowski/easymotion/service/grpcservice/proto/opencv"
	"github.com/rlaskowski/manage"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedVideoServer
	listen  net.Listener
	grpcSrv *grpc.Server
	address string
	context context.Context
	cancel  context.CancelFunc
}

func (GrpcServer) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.server.grpc",
		Priority:  2,
		Intstance: newServer(config.DefaultOptions.GrpcAddress),
	}
}

func newServer(address string) *GrpcServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &GrpcServer{
		grpcSrv: grpc.NewServer(),
		address: address,
		context: ctx,
		cancel:  cancel,
	}
}

func (g *GrpcServer) Start() error {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("grpc server start error: %s", err.Error())
	}

	pb.RegisterVideoServer(g.grpcSrv, g)

	return g.grpcSrv.Serve(listen)
}

func (g *GrpcServer) Stop() error {
	return nil
}

func (g *GrpcServer) Stream(req *pb.VideoRequest, resp pb.Video_StreamServer) error {
	return nil
}

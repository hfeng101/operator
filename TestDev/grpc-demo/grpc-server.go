package grpc_demo

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/grpc"
	grpc "google.golang.org/grpc"
	"net"
	pb "path/to/you"
)

type TestGRPCServer struct {
	pb.UnimplementedTestGRPCServer
}

func (*TestGRPCServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.GetA() + req.GetB()
	return &pb.AddResponse{Result: result}, nil
}

func StartGRPCServer() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Listen for localhost:12345 failed, err is %s", err.Error())
	}

	grpcHandle := iris.New()
	rootParty := grpcHandle.Party("/")
	rootParty.HandleFunc("Post", "/", func() {
		fmt.Println("Welcome to grpc server ,haha !")
	})
	grpcServer := grpc.New(grpcHandle)

	pb.RegisterTestGRPCServer(grpcServer, &TestGRPCServer{})

	fmt.Println("grpc server is running on port localhost:12345")

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Println("grpc server serve failed, err is %s", err)
	}
}

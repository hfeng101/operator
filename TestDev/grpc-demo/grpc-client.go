package grpc_demo

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func startGrpcClient() {
	//listener, err := net.Listen("tcp", ":12345")
	//if err != nil {
	//	fmt.Println("listen tcp dial failed, err is %s", err.Error())
	//	return
	//}

	conn, err := grpc.Dial("127.0.0.1:12345")
	if err != nil {
		fmt.Println("Connet to grpc server failed, err is %s", err.Error())
		return
	}

	defer conn.Close()

	client := pb.NewTestGRPCClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel

	req := pb.AddRequest{A: 10, B: 20}
	resp, err := client.Add(ctx, req)
	if err != nil {
		fmt.Println("grpc add function failed, err is %s", err.Error())
	}

	fmt.Println("call add function by grpc, a:%d,b:%d,result is %d", req.A, req.B, resp.Result)
}

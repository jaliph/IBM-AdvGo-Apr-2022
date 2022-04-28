package main

import (
	"context"
	"errors"
	"fmt"
	"grpc-demo/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type appServiceImpl struct {
	proto.UnimplementedAppServiceServer
}

func (asi *appServiceImpl) Add(ctx context.Context, req *proto.AddRequest) (res *proto.AddResponse, er error) {
	x := req.GetX()
	y := req.GetY()
	fmt.Printf("Add Operation invoked with x=%d and y=%d\n", x, y)
	timeOut := time.After(10 * time.Second)

LOOP:
	for {
		select {
		case <-timeOut:
			result := x + y
			res = &proto.AddResponse{
				Result: result,
			}
			break LOOP
		case <-ctx.Done():
			fmt.Println("Cancel instruction received")
			er = errors.New("interrupt received")
			break LOOP
		}
	}
	return
}

func (asi *appServiceImpl) GeneratePrimes(req *proto.PrimeRequest, serverStream proto.AppService_GeneratePrimesServer) error {
	start := req.GetStart()
	end := req.GetEnd()
	for no := start; no <= end; no++ {
		if isPrime(no) {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Generated prime no : %d\n", no)
			res := &proto.PrimeResponse{
				PrimeNo: no,
			}
			serverStream.Send(res)
		}
	}
	return nil
}

func isPrime(no int32) bool {
	for i := int32(2); i <= (no / 2); i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	asi := &appServiceImpl{}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAppServiceServer(grpcServer, asi)
	grpcServer.Serve(listener)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dvaganov/fiboser/pkg/delivery/fibogrpc"
	"google.golang.org/grpc"
)

type (
	grpcServer struct {
		fiboSrv fibogrpc.FibonacciServer
	}
)

func (s *grpcServer) Run(ctx context.Context, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer()
	fibogrpc.RegisterFibonacciServer(grpcSrv, s.fiboSrv)

	go func() {
		<-ctx.Done()
		grpcSrv.Stop()
		log.Println(lis.Close())
	}()

	go func() {
		if err := grpcSrv.Serve(lis); err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	return nil
}

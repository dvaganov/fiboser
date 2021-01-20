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
		srv     *grpc.Server
		lis     net.Listener
	}
)

func (s *grpcServer) Run(port int) error {
	var err error
	s.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	fibogrpc.RegisterFibonacciServer(s.srv, s.fiboSrv)

	go func() {
		if err := s.srv.Serve(s.lis); err != grpc.ErrServerStopped {
			log.Println(err)
		}
	}()

	return nil
}

func (s *grpcServer) Stop(ctx context.Context) error {
	s.srv.Stop()
	return s.lis.Close()
}

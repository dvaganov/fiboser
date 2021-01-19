package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dvaganov/fiboser/pkg/delivery/fibogrpc"
	"github.com/dvaganov/fiboser/pkg/delivery/rest"
	"github.com/dvaganov/fiboser/pkg/fibonacci"
)

var port = flag.Int("port", 8080, "the port for rest server")
var grpcPort = flag.Int("grpcPort", 8081, "the port for grpc server")

func main() {
	flag.Parse()
	service := fibonacci.NewService()
	ctrl := rest.NewController(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", ctrl.GetFibonacciRange)

	srv := http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", *port),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
		srv.Shutdown(ctx)
	}()

	grpcSrv := grpcServer{
		fiboSrv: fibogrpc.NewFiboServer(service),
	}
	err := grpcSrv.Run(ctx, *grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatalln(srv.ListenAndServe())
}

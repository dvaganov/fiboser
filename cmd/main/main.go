package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dvaganov/fiboser/pkg/delivery/fibogrpc"
	"github.com/dvaganov/fiboser/pkg/delivery/fiborest"
	"github.com/dvaganov/fiboser/pkg/fibonacci"
)

var restPort = flag.Int("restPort", 8080, "the port for rest server")
var grpcPort = flag.Int("grpcPort", 8081, "the port for grpc server")

const (
	shutDownTimeout = 5 * time.Second
)

func main() {
	flag.Parse()
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	service := fibonacci.NewService()

	restSrv := restServer{fiboCtrl: fiborest.NewController(service)}
	err := restSrv.Run(*restPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcServer{fiboSrv: fibogrpc.NewFiboServer(service)}
	err = grpcSrv.Run(*grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	<-osSignal
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
		defer cancel()

		if err := restSrv.Stop(ctx); err != nil {
			log.Println(err)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
		defer cancel()

		if err := grpcSrv.Stop(ctx); err != nil {
			log.Println(err)
		}

		wg.Done()
	}()

	wg.Wait()
}

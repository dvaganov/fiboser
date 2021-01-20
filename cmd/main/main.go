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
var redisDsn = flag.String("redisDsn", "", "data source name for connection like redis://localhost:6379/?db=0&password=")

const (
	shutDownTimeout = 5 * time.Second
)

func main() {
	flag.Parse()
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var cache fibonacci.Cache

	if *redisDsn != "" {
		rdb, err := NewRedis(*redisDsn)
		if err != nil {
			log.Fatal(err)
		}
		cache = fibonacci.NewRedisCache(rdb)
	}

	service := fibonacci.NewService(cache)

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

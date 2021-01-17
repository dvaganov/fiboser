package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dvaganov/fiboser/pkg/delivery/rest"
	"github.com/dvaganov/fiboser/pkg/fibonacci"
)

var port = flag.Int("port", 8080, "the port for rest server")

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

	log.Fatalln(srv.ListenAndServe())
}
